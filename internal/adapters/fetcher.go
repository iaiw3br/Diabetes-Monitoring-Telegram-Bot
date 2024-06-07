package adapters

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"notifier/internal/entities"
)

type HttpFetcher struct{}

func NewHttpFetcher() *HttpFetcher {
	return &HttpFetcher{}
}

func (hf *HttpFetcher) FetchEntry() (*entities.SGVResponse, error) {
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

	var response []entities.SGVResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response[0], nil
}

func (hf *HttpFetcher) FetchTreatments() (*entities.TreatmentResponse, error) {
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

	var response []entities.TreatmentResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response[0], nil
}
