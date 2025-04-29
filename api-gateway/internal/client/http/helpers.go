package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func postJSON(url string, body interface{}) (map[string]interface{}, error) {
	jsonData, _ := json.Marshal(body)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	return decodeBody(resp)
}

func getJSON(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	return decodeBody(resp)
}

func putJSON(url string, body interface{}) (map[string]interface{}, error) {
	jsonData, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	return decodeBody(resp)
}

func deleteRequest(url string) error {
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode >= 400 {
		return errors.New("error from backend")
	}
	return nil
}

func decodeBody(resp *http.Response) (map[string]interface{}, error) {
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(body))
	}

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
