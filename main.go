package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// Metadata структура для metadata.json
type Metadata struct {
	Flavor string `json:"flavor"`

	Fields []struct {
		ID    string      `json:"id"`
		Value interface{} `json:"value"`
	} `json:"fields"`
}

func main() {
	metadataData := []Metadata{
		{
			Flavor: "dublincore/episode",
			Fields: []struct {
				ID    string      `json:"id"`
				Value interface{} `json:"value"`
			}{
				{
					ID:    "title",
					Value: "gotest",
				},
				{
					ID:    "subjects",
					Value: []string{"John Clark", "Thiago Melo Costa"},
				},
				{
					ID:    "description",
					Value: "A great description",
				},
				{
					ID:    "startDate",
					Value: "2016-06-22",
				},
				{
					ID:    "startTime",
					Value: "13:30:00Z",
				},
			},
		},
	}

	metadataJSON, err := json.Marshal(metadataData)
	if err != nil {
		fmt.Println("Error marshaling metadata data:", err)
		return
	}

	presenterData, err := os.ReadFile("C:\\Users\\Zanzhit\\opencast\\test.mkv")
	if err != nil {
		fmt.Println("Error reading presenter file:", err)
		return
	}

	aclJSON, err := os.ReadFile("acl.json")
	if err != nil {
		fmt.Println("Error reading acl file:", err)
		return
	}

	processingJSON, err := os.ReadFile("process.json")
	if err != nil {
		fmt.Println("Error reading processing file:", err)
		return
	}

	// Формирование данных в формате multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавление JSON к multipart форме
	files := map[string][]byte{
		"presenter":  presenterData,
		"acl":        aclJSON,
		"metadata":   metadataJSON,
		"processing": processingJSON,
	}

	// Добавление файлов к multipart форме
	for fieldName, fileData := range files {
		if fieldName == "presentation" || fieldName == "presenter" {
			// Добавление видеофайла как файла
			part, err := writer.CreateFormFile(fieldName, fieldName+".mkv")
			if err != nil {
				fmt.Println("Error creating form file:", err)
				return
			}
			_, err = io.Copy(part, bytes.NewReader(fileData))
			if err != nil {
				fmt.Println("Error copying file data:", err)
				return
			}
		} else {
			// Добавление текстовых данных в форму
			part, err := writer.CreateFormField(fieldName)
			if err != nil {
				fmt.Println("Error creating form field:", err)
				return
			}
			part.Write(fileData)
		}
	}

	// Завершение формы
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return
	}

	// Формирование запроса
	req, err := http.NewRequest("POST", "http://localhost:8080/api/events", body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Установка заголовков запроса
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.SetBasicAuth("admin", "opencast")

	// Отправка запроса
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Чтение ответа
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// // Запись ответа в файл
	// err = os.WriteFile("resp.json", respBody, 0644)
	// if err != nil {
	// 	fmt.Println("Error writing response to file:", err)
	// 	return
	// }

	fmt.Println(string(respBody))
}
