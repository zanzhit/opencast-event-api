package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	// Чтение содержимого файлов для передачи в запрос
	presenterData, err := os.ReadFile("test.mkv")
	if err != nil {
		fmt.Println("Error reading presenter file:", err)
		return
	}

	aclData, err := os.ReadFile("acl.json")
	if err != nil {
		fmt.Println("Error reading acl file:", err)
		return
	}

	metadataData, err := os.ReadFile("metadata.json")
	if err != nil {
		fmt.Println("Error reading metadata file:", err)
		return
	}

	processingData, err := os.ReadFile("process.json")
	if err != nil {
		fmt.Println("Error reading processing file:", err)
		return
	}

	// Формирование данных в формате multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавление файлов к multipart форме
	files := map[string][]byte{
		"presenter":  presenterData,
		"acl":        aclData,
		"metadata":   metadataData,
		"processing": processingData,
	}

	for fieldName, fileData := range files {
		part, err := writer.CreateFormField(fieldName)
		if err != nil {
			fmt.Println("Error creating form field:", err)
			return
		}
		part.Write(fileData)
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

	// Запись ответа в файл
	err = os.WriteFile("resp.json", respBody, 0644)
	if err != nil {
		fmt.Println("Error writing response to file:", err)
		return
	}

	fmt.Println("Response written to resp.json")
}
