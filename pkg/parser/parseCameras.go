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

func ParseCameras(pathToShinobiCameras string, camerasJSON []byte) (metadatas [][]byte, presenters [][]byte, shinobiFiles []string, err error) {
	var cameras []Camera

	err = json.Unmarshal(camerasJSON, &cameras)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error unmarshal: %v", err)
	}

	for _, cam := range cameras {
		path := fmt.Sprintf("%s/%s", pathToShinobiCameras, cam.ID)
		files, err := os.ReadDir(path)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error reading directory: %v", err)
		}

		for _, file := range files {
			fileWithPath := fmt.Sprintf("%s/%s", path, file.Name())
			shinobiFiles = append(shinobiFiles, fileWithPath)

			presenter, err := os.ReadFile(fileWithPath)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("error reading presenter file: %v", err)
			}
			presenters = append(presenters, presenter)

			startDate := file.Name()[:10]
			startTime := strings.ReplaceAll(file.Name(), "-", ":")[11 : len(file.Name())-4]

			duration, err := videoDuration(fileWithPath)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("error videoDuration: %v", err)
			}

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
						// {
						// 	ID:    "duration",
						// 	Value: "00:00:20",
						// },
						{
							ID:    "location",
							Value: cam.IP,
						},
					},
				},
			}

			metadataJSON, err := json.Marshal(metadata)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("error marshal: %v", err)
			}

			metadatas = append(metadatas, metadataJSON)
		}
	}

	return metadatas, presenters, shinobiFiles, nil
}
