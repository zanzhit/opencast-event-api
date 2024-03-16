package request

import (
	"bytes"
	"encoding/json"

	"opencast/pkg/parser"
)

func SendMultipartRequest(cameraData parser.VideoData, acl, processing []byte, configWithPath []byte) ([]byte, error) {
	var config config

	err := json.Unmarshal(configWithPath, &config)
	if err != nil {
		return nil, err
	}

	body := &bytes.Buffer{}
	files := map[string][]byte{
		"presenter":  cameraData.Presenter,
		"acl":        acl,
		"metadata":   cameraData.Metadata,
		"processing": processing,
	}

	header, err := AddMultipartForm(files, body)
	if err != nil {
		return nil, err
	}

	req, err := FormingNewRequest(body, header, config.Login, config.Password, config.OpencastURL)
	if err != nil {
		return nil, err
	}

	respBody, err := SendRequest(req)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
