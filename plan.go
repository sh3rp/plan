package plan

type Version struct {
	ApiVersion    string `json:"api_ver"`
	ServerVersion string `json:"server_ver"`
}

type Plan struct {
	Id         string       `json:"id"`
	Body       string       `json:"body"`
	PostedTime int64        `json:"timestamp"`
	Links      []string     `json:"links"`
	Tags       []string     `json:"tags"`
	Location   *GeoLocation `json:"location"`
}

type PlanInfo struct {
	Handle    string `json:"handle"`
	RealName  string `json:"real_name"`
	Location  string `json:"location"`
	Homepage  string `json:"homepage"`
	AvatarURL string `json:"avatar_url"`
}

type GeoLocation struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
}
