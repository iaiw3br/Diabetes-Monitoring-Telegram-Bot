package entities

import (
	"time"
)

type SGVResponse struct {
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
	ID                string    `json:"_id"`
	Timestamp         int       `json:"timestamp"`
	EventType         string    `json:"eventType"`
	EnteredBy         string    `json:"enteredBy"`
	UID               string    `json:"uuid"`
	Insulin           float64   `json:"insulin"`
	InsulinInjections string    `json:"insulinInjections"`
	CreatedAt         time.Time `json:"created_at"`
	SysTime           string    `json:"sysTime"`
	UtcOffset         int       `json:"utcOffset"`
	Mills             int64     `json:"mills"`
	Carbs             *int      `json:"carbs"`
}
