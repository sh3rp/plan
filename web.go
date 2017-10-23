package plan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type WebService struct {
	PlanDB     PlanDB
	authString string
}

type Response struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type PlanResponse struct {
	Response
	Plans []*Plan `json:"plans"`
}

type PlanInfoResponse struct {
	Response
	PlanInfo *PlanInfo `json:"info"`
}

func (ws *WebService) Start(port int, authString string) {
	ws.authString = authString
	http.HandleFunc("/now", ws.current)
	http.HandleFunc("/info", ws.info)
	http.HandleFunc("/plan/", ws.plan)
	http.HandleFunc("/plans", ws.plans)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (ws *WebService) current(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ws.getCurrent(w, r)
	case "POST":
		fallthrough
	case "PUT":
		ws.postNew(w, r)
	}
}

func (ws *WebService) info(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ws.currentInfo(w, r)
	case "POST":
		fallthrough
	case "PUT":
		ws.postInfo(w, r)
	}
}

func (ws *WebService) plan(w http.ResponseWriter, r *http.Request) {
	tokens := strings.Split(r.URL.Path, "/")
	id := tokens[len(tokens)-1]

	plan, err := ws.PlanDB.GetPlan(id)

	if err != nil {
		sendError(w, 101, "Error retrieving plan")
	} else {
		sendPlan(w, plan)
	}
}

func (ws *WebService) plans(w http.ResponseWriter, r *http.Request) {
	plans, err := ws.PlanDB.GetPlans()

	if err != nil {
		sendError(w, 101, "Error retrieving plans")
	} else {
		sendPlans(w, plans)
	}
}

func (ws *WebService) getCurrent(w http.ResponseWriter, r *http.Request) {
	plan, err := ws.PlanDB.GetCurrentPlan()

	if err != nil {
		sendError(w, 101, "Error getting current plan")
	} else {
		sendPlan(w, plan)
	}
}

func (ws *WebService) postNew(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("x-plan-auth") != ws.authString {
		sendError(w, 102, "Authentication error, wrong password")
		return
	}
	var plan Plan
	err := json.NewDecoder(r.Body).Decode(&plan)
	if err != nil {
		sendError(w, 104, fmt.Sprintf("Error decoding plan object: %v", err))
		return
	}
	newPlan, err := ws.PlanDB.NewPlan(&plan)

	if err != nil {
		sendError(w, 105, "Error creating new plan")
	} else {
		sendPlan(w, newPlan)
	}
}

func (ws *WebService) currentInfo(w http.ResponseWriter, r *http.Request) {
	info := ws.PlanDB.GetInfo()
	sendPlanInfo(w, info)
}

func (ws *WebService) postInfo(w http.ResponseWriter, r *http.Request) {
	var info PlanInfo
	json.NewDecoder(r.Body).Decode(&info)
	ws.PlanDB.SetInfo(&info)
	sendPlanInfo(w, &info)
}

func sendError(w http.ResponseWriter, code int, msg string) {
	json.NewEncoder(w).Encode(&PlanResponse{
		Response: Response{
			Code:      code,
			Message:   msg,
			Timestamp: timestamp(),
		},
		Plans: []*Plan{},
	})
}

func sendPlan(w http.ResponseWriter, plan *Plan) {
	var plans []*Plan
	if plan == nil {
		plans = []*Plan{}
	} else {
		plans = []*Plan{plan}
	}
	json.NewEncoder(w).Encode(&PlanResponse{
		Response: Response{
			Code:      1,
			Message:   "Ok",
			Timestamp: timestamp(),
		},
		Plans: plans,
	})
}

func sendPlans(w http.ResponseWriter, plans []*Plan) {
	json.NewEncoder(w).Encode(&PlanResponse{
		Response: Response{
			Code:      1,
			Message:   "Ok",
			Timestamp: timestamp(),
		},
		Plans: plans,
	})
}

func sendPlanInfo(w http.ResponseWriter, info *PlanInfo) {
	json.NewEncoder(w).Encode(&PlanInfoResponse{
		Response: Response{
			Code:      1,
			Message:   "Ok",
			Timestamp: timestamp(),
		},
		PlanInfo: info,
	})
}

func timestamp() int64 {
	return int64(time.Now().UnixNano() / 1000000)
}
