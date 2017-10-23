package plan

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/oklog/ulid"
)

type PlanDB interface {
	NewPlan(*Plan) (*Plan, error)
	SavePlan(*Plan) (*Plan, error)
	DeletePlan(string) error
	GetCurrentPlan() (*Plan, error)
	GetPlan(string) (*Plan, error)
	GetPlans() ([]*Plan, error)

	SetInfo(*PlanInfo) error
	GetInfo() *PlanInfo
}

var PLAN_BUCKET = []byte("plan")
var CURRENT_KEY = []byte("CURRENT")
var INFO_KEY = []byte("INFO")

type BoltPlanDB struct {
	db *bolt.DB
}

func NewBoltPlanDB(dbdir string) (PlanDB, error) {
	db, err := bolt.Open(dbdir+"/plan.db", 0755, nil)

	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(PLAN_BUCKET)
		return err
	})

	return &BoltPlanDB{
		db: db,
	}, err
}

func (b *BoltPlanDB) SetInfo(pi *PlanInfo) error {
	bytes := &bytes.Buffer{}
	err := json.NewEncoder(bytes).Encode(pi)

	if err != nil {
		return err
	}

	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(PLAN_BUCKET)

		return bucket.Put(INFO_KEY, bytes.Bytes())
	})
}

func (b *BoltPlanDB) GetInfo() *PlanInfo {
	bytes := &bytes.Buffer{}
	b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(PLAN_BUCKET)
		bytes.Write(bucket.Get(INFO_KEY))
		return nil
	})

	var info PlanInfo

	json.NewDecoder(bytes).Decode(&info)

	return &info
}

func (b *BoltPlanDB) NewPlan(p *Plan) (*Plan, error) {
	p.Id = newId()
	if p.Links == nil {
		p.Links = []string{}
	}
	if p.Tags == nil {
		p.Tags = []string{}
	}
	if p.Location == nil {
		p.Location = &GeoLocation{}
	}
	p.PostedTime = int64(time.Duration(time.Now().UnixNano()) / time.Millisecond)
	plan, err := b.SavePlan(p)

	if err != nil {
		return nil, err
	}

	return plan, nil
}

func (b *BoltPlanDB) SavePlan(p *Plan) (*Plan, error) {
	err := b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(PLAN_BUCKET)
		return bucket.Put(idToBytes(p.Id), ToJSON(p))
	})

	return p, err
}

func (b *BoltPlanDB) DeletePlan(id string) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(PLAN_BUCKET)
		return bucket.Delete(idToBytes(id))
	})
	return err
}

func (b *BoltPlanDB) GetCurrentPlan() (*Plan, error) {
	var plan *Plan
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(PLAN_BUCKET)

		cursor := bucket.Cursor()
		prefix := []byte("/p/")
		for k, v := cursor.Last(); k != nil; k, v = cursor.Prev() {
			if bytes.HasPrefix(k, prefix) {
				plan = ToOBJ(v)
				break
			}
		}
		return nil
	})
	return plan, err
}

func (b *BoltPlanDB) GetPlan(id string) (*Plan, error) {
	var plan *Plan
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(PLAN_BUCKET)
		bytes := bucket.Get(idToBytes(id))
		plan = ToOBJ(bytes)
		return nil
	})
	return plan, err
}

func (b *BoltPlanDB) GetPlans() ([]*Plan, error) {
	var plans []*Plan
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(PLAN_BUCKET)
		cursor := bucket.Cursor()

		prefix := []byte("/p/")
		for k, v := cursor.Last(); k != nil; k, v = cursor.Prev() {
			if bytes.HasPrefix(k, prefix) {
				plans = append(plans, ToOBJ(v))
			}
		}
		return nil
	})

	if plans == nil {
		plans = []*Plan{}
	}

	return plans, err
}

func newId() string {
	id := ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader)

	return strings.ToLower(id.String())
}

func idToBytes(id string) []byte {
	buf := &bytes.Buffer{}
	buf.Write([]byte("/p/"))
	buf.Write([]byte(id))
	return buf.Bytes()
}

func ToJSON(p *Plan) []byte {
	bytes := &bytes.Buffer{}
	json.NewEncoder(bytes).Encode(p)
	return bytes.Bytes()
}

func ToOBJ(data []byte) *Plan {
	plan := &Plan{}
	bytes := &bytes.Buffer{}
	bytes.Write(data)
	json.NewDecoder(bytes).Decode(plan)
	return plan
}
