package request

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func AddMultipartForm(files map[string][]byte, body *bytes.Buffer) (string, error) {
	writer := multipart.NewWriter(body)
	defer writer.Close()

	for fieldName, fileData := range files {
		if fieldName == "presentation" || fieldName == "presenter" || fieldName == "audio" {
			part, err := writer.CreateFormFile(fieldName, fieldName+".mp4")
			if err != nil {
				return "", err
			}

			_, err = io.Copy(part, bytes.NewReader(fileData))
			if err != nil {
				return "", err
			}

			continue
		}

		part, err := writer.CreateFormField(fieldName)
		if err != nil {
			return "", err
		}
		part.Write(fileData)
	}

	return writer.FormDataContentType(), nil
}

func FormingNewRequest(body *bytes.Buffer, value string, login, password, url string) (*http.Request, error) {
	url = fmt.Sprintf("%s/api/events", url)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", value)
	req.SetBasicAuth(login, password)

	return req, nil
}

func SendRequest(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
