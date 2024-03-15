package request

type config struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	OpencastURL string `json:"urlOpencast"`
	APIkey      string `json:"apiKey"`
	GroupKey    string `json:"groupKey"`
	ShinobiURL  string `json:"urlShinobi"`
}
