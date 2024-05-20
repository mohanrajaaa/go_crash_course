package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

const serverURL = "http://localhost:8080"

type File struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func CreateFile(filename, content string) error {
	filePath := filepath.Join(serverURL, "files")
	fileData := File{Name: filename, Content: content}
	jsonData, err := json.Marshal(fileData)
	if err != nil {
		return err
	}

	resp, err := http.Post(filePath, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create file, status code: %d", resp.StatusCode)
	}

	return nil
}

func ReadFile(filename string) (string, error) {
	filePath := filepath.Join(serverURL, "files", filename)
	resp, err := http.Get(filePath)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to read file, status code: %d", resp.StatusCode)
	}

	var fileContent string
	json.NewDecoder(resp.Body).Decode(&fileContent)
	return fileContent, nil
}

func UpdateFile(filename, content string) error {
	filePath := filepath.Join(serverURL, "files", filename)
	fileData := File{Content: content}
	jsonData, err := json.Marshal(fileData)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, filePath, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update file, status code: %d", resp.StatusCode)
	}

	return nil
}

func DeleteFile(filename string) error {
	filePath := filepath.Join(serverURL, "files", filename)
	req, err := http.NewRequest(http.MethodDelete, filePath, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete file, status code: %d", resp.StatusCode)
	}

	return nil
}
