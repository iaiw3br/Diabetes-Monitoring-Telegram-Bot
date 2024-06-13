package usecases

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestChecker_IsInterval(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		lastSent       time.Time
		interval       time.Duration
		waitDuration   time.Duration
		expectedResult bool
	}{
		{
			name:           "Interval has passed",
			lastSent:       time.Now().Add(-10 * time.Second),
			interval:       5 * time.Second,
			waitDuration:   0,
			expectedResult: true,
		},
		{
			name:           "Interval has not passed",
			lastSent:       time.Now(),
			interval:       5 * time.Second,
			waitDuration:   0,
			expectedResult: false,
		},
		{
			name:           "Interval passes after wait",
			lastSent:       time.Now(),
			interval:       2 * time.Second,
			waitDuration:   3 * time.Second,
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			c := &Checker{
				lastSent: tt.lastSent,
				interval: tt.interval,
			}

			if tt.waitDuration > 0 {
				time.Sleep(tt.waitDuration)
			}

			// Act
			result := c.isInterval()

			// Assert
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestIsLongTimeAgo(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		date           time.Time
		localTime      time.Time
		difference     time.Duration
		expectedResult bool
	}{
		{
			name:           "Time difference is greater than specified duration",
			date:           time.Now().Add(-10 * time.Minute),
			localTime:      time.Now(),
			difference:     5 * time.Minute,
			expectedResult: true,
		},
		{
			name:           "Time difference is less than specified duration",
			date:           time.Now().Add(-2 * time.Minute),
			localTime:      time.Now(),
			difference:     5 * time.Minute,
			expectedResult: false,
		},
		{
			name:           "Time difference is exactly the specified duration",
			date:           time.Now().Add(-5 * time.Minute),
			localTime:      time.Now(),
			difference:     5 * time.Minute,
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := isLongTimeAgo(tt.date, tt.localTime, tt.difference)

			// Assert
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestGetMessage(t *testing.T) {
	tests := []struct {
		name           string
		prevSGV        float64
		value          float64
		expectedResult string
	}{
		{
			name:           "Value above highThreshold",
			value:          12.0,
			expectedResult: fmt.Sprintf(highText, 12.0),
		},
		{
			name:           "Value below lowThreshold",
			value:          4.0,
			expectedResult: fmt.Sprintf(lowText, 4.0),
		},
		{
			name:           "Value within thresholds, previous and current within thresholds",
			prevSGV:        8.0,
			value:          9.0,
			expectedResult: "",
		},
		{
			name:           "Value within thresholds, previous above highThreshold",
			prevSGV:        11.0,
			value:          9.0,
			expectedResult: "",
		},
		{
			name:           "Value within thresholds, previous below lowThreshold",
			prevSGV:        4.0,
			value:          5.0,
			expectedResult: "",
		},
		{
			name:           "Value within thresholds, previous below lowThreshold",
			prevSGV:        4.6,
			value:          4.4,
			expectedResult: fmt.Sprintf(lowText, 4.4),
		},
		{
			name:           "Value within thresholds, previous below lowThreshold",
			prevSGV:        9.9,
			value:          10.1,
			expectedResult: fmt.Sprintf(highText, 10.1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Checker{
				prevSGV: tt.prevSGV,
			}
			result := c.getMessage(tt.value)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestCalculateValue(t *testing.T) {
	tests := []struct {
		name           string
		sgv            int
		expectedResult float64
	}{
		{
			name:           "Round value",
			sgv:            180,
			expectedResult: 10.0,
		},
		{
			name:           "Round down",
			sgv:            175,
			expectedResult: 9.7,
		},
		{
			name:           "Round up",
			sgv:            185,
			expectedResult: 10.3,
		},
		{
			name:           "Round to nearest tenth",
			sgv:            190,
			expectedResult: 10.6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateValue(tt.sgv)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
