package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
)

func GETRequest[T any](url string) (T, error) {
	var data T
	resp, err := http.Get(url)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	if resp.StatusCode != http.StatusOK {
		if body != nil {
			err = json.Unmarshal(body, &data)
			if err != nil {
				return data, err
			}
			return data, errors.New("GET request failed with status code: " + resp.Status)
		}
		return data, errors.New("GET request failed with status code: " + resp.Status)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func POSTRequest[T any](endpoint string, data *fiber.Map) (T, error) {
	var responseData T
	urlPost := fmt.Sprintf("%s%s%s", ServerUrl, "/api/v1/", endpoint)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return responseData, err
	}

	req, err := http.NewRequest("POST", urlPost, bytes.NewBuffer(jsonData))
	if err != nil {
		return responseData, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return responseData, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return responseData, err
	}

	if resp.StatusCode != http.StatusOK {
		if body != nil {
			err = json.Unmarshal(body, &responseData)
			if err != nil {
				return responseData, err
			}
			return responseData, errors.New("GET request failed with status code: " + resp.Status)
		}
		return responseData, errors.New("GET request failed with status code: " + resp.Status)
	}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return responseData, err
	}

	return responseData, nil
}
