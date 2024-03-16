package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Metadata struct {
	Flavor string `json:"flavor"`

	Fields []struct {
		ID    string      `json:"id"`
		Value interface{} `json:"value"`
	} `json:"fields"`
}

func ParseCameras(pathToShinobiCameras string, camerasJSON []byte) ([]VideoData, error) {
	var cameras []Camera
	var dataFromVideos []VideoData

	err := json.Unmarshal(camerasJSON, &cameras)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal: %v", err)
	}

	for _, cam := range cameras {
		path := fmt.Sprintf("%s/%s", pathToShinobiCameras, cam.ID)

		files, err := os.ReadDir(path)
		if err != nil {
			return nil, fmt.Errorf("error reading directory: %v", err)
		}

		for _, file := range files {
			fileWithPath := fmt.Sprintf("%s/%s", path, file.Name())

			duration, err := videoDuration(fileWithPath)
			if err != nil {
				fmt.Printf("issues with %s : %s", file.Name(), err)
				os.Remove(fileWithPath)
				continue
			}

			presenter, err := os.ReadFile(fileWithPath)
			if err != nil {
				return nil, fmt.Errorf("error reading presenter file: %v", err)
			}

			startDate := file.Name()[:10]
			startTime := strings.ReplaceAll(file.Name(), "-", ":")[11 : len(file.Name())-4]

			metadata := []Metadata{
				{
					Flavor: "dublincore/episode",
					Fields: []struct {
						ID    string      `json:"id"`
						Value interface{} `json:"value"`
					}{
						{
							ID:    "title",
							Value: file.Name(),
						},
						{
							ID:    "startDate",
							Value: startDate,
						},
						{
							ID:    "startTime",
							Value: startTime,
						},
						{
							ID:    "duration",
							Value: duration,
						},
						{
							ID:    "location",
							Value: cam.IP,
						},
					},
				},
			}

			metadataJSON, err := json.Marshal(metadata)
			if err != nil {
				return nil, fmt.Errorf("error marshal: %v", err)
			}

			dataFromVideos = append(dataFromVideos, VideoData{Metadata: metadataJSON, Presenter: presenter, ShinobiFile: fileWithPath})
		}
	}

	return dataFromVideos, nil
}
