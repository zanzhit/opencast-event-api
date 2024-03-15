package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func DeleteVideo(conf []byte, video string) (string, error) {
	var config config

	err := json.Unmarshal(conf, &config)
	if err != nil {
		return "", err
	}

	shinobiReq := fmt.Sprintf("%s/%s/%s/%s", config.ShinobiURL, config.APIkey, video[strings.Index(video, "videos"):], "delete")

	client := &http.Client{}

	req, err := http.NewRequest("GET", shinobiReq, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return resp.Status, nil
}
