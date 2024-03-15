package request

import (
	"bytes"
	"encoding/json"
)

func SendMultipartRequest(presenter, acl, metadata, processing []byte, configWithPath []byte) ([]byte, error) {
	var config config

	err := json.Unmarshal(configWithPath, &config)
	if err != nil {
		return nil, err
	}

	body := &bytes.Buffer{}
	files := map[string][]byte{
		"presenter":  presenter,
		"acl":        acl,
		"metadata":   metadata,
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
