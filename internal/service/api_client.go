package service

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"notifier/pkg/models"
)

func FetchData() (*models.Response, error) {
	requestURL := "https://j82719866.nightscout-jino.ru/api/v1/entries.json?count=1"
	res, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response []models.Response
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response[0], nil
}

func FetchTreatments() (*models.TreatmentResponse, error) {
	requestURL := "https://j82719866.nightscout-jino.ru/api/v1/treatments?count=1"
	res, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response []models.TreatmentResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response[0], nil
}
