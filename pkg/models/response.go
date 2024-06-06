package models

import (
	"github.com/google/uuid"
	"time"
)

type Response struct {
	Device     string    `json:"device"`
	Date       int       `json:"date"`
	DateString time.Time `json:"dateString"`
	SGV        int       `json:"sgv"`
	Delta      float64   `json:"delta"`
	Direction  string    `json:"direction"`
	Type       string    `json:"type"`
	Filtered   float64   `json:"filtered"`
	Unfiltered float64   `json:"unfiltered"`
	RSSI       int       `json:"rssi"`
	Noise      int       `json:"noise"`
	SysTime    time.Time `json:"sysTime"`
	UtcOffset  int       `json:"utcOffset"`
	ID         string    `json:"_id"`
	Mills      int       `json:"mills"`
}

type TreatmentResponse struct {
	UID       uuid.UUID `json:"id"`
	CreatedAt string    `json:"created_at"`
	Insulin   float64   `json:"insulin"`
}
