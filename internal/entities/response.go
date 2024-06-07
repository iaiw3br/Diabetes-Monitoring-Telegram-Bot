package entities

import (
	"github.com/google/uuid"
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
	Timestamp int `json:"timestamp"`
	// The type of treatment event
	EventType string `json:"eventType"`
	// Who entered the treatment.
	EnteredBy string    `json:"enteredBy"`
	UID       uuid.UUID `json:"uuid"`
	// Amount of carbs consumed in grams
	Carbs             int    `json:"carbs"`
	InsulinInjections string `json:"insulinInjections"`
	// The date of the event, might be set retroactively
	CreatedAt time.Time `json:"created_at"`
	SysTime   time.Time `json:"sysTime"`
	// Internally assigned id
	ID        string `json:"_id"`
	UtcOffset int    `json:"utcOffset"`
	Mills     int    `json:"mills"`
	// Amount of insulin, if any.
	Insulin float64 `json:"insulin"`
}
