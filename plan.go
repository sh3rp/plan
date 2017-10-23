package plan

type Plan struct {
	Id         string   `json:"id"`
	Body       string   `json:"body"`
	PostedTime int64    `json:"timestamp"`
	Links      []string `json:"links"`
	Tags       []string `json:"tags"`
}

type PlanInfo struct {
	Handle    string `json:"handle"`
	RealName  string `json:"real_name"`
	Location  string `json:"location"`
	Homepage  string `json:"homepage"`
	AvatarURL string `json:"avatar_url"`
}
