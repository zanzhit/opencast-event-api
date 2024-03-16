package parser

type Camera struct {
	IP string `json:"ip"`
	ID string `json:"id"`
}

type VideoData struct {
	Metadata    []byte
	Presenter   []byte
	ShinobiFile string
}
