package usecases

import (
	"testing"
	"time"
)

func Test_formatDifference(t *testing.T) {
	type args struct {
		oldSGV     float64
		currentSGV float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Positive difference",
			args: args{oldSGV: 120.0, currentSGV: 100.0},
			want: "+20.0",
		},
		{
			name: "Negative difference",
			args: args{oldSGV: 100.0, currentSGV: 120.0},
			want: "-20.0",
		},
		{
			name: "Zero difference",
			args: args{oldSGV: 100.0, currentSGV: 100.0},
			want: "+0.0",
		},
		{
			name: "Positive floating point difference",
			args: args{oldSGV: 100.0, currentSGV: 99.5},
			want: "+0.5",
		},
		{
			name: "Negative floating point difference",
			args: args{oldSGV: 99.5, currentSGV: 100},
			want: "-0.5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatDifference(tt.args.oldSGV, tt.args.currentSGV); got != tt.want {
				t.Errorf("formatDifference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stopTimerIfExists(t *testing.T) {
	t.Run("Timer exists and is stopped", func(t *testing.T) {
		timerMap := make(map[string]*time.Timer)

		// Create a timer and add it to the map
		timer := time.AfterFunc(10*time.Second, func() {})
		key := "test_timer"
		timerMap[key] = timer

		// Ensure the timer is running
		if !timer.Stop() {
			t.Fatalf("Expected timer to be running")
		}

		// Start the timer again for the test
		timer.Reset(10 * time.Second)

		// Call the function to stop the timer
		stopTimerIfExists(timerMap, key)

		// Check if the timer is stopped
		if timer.Stop() {
			t.Errorf("Expected timer to be stopped, but it is still running")
		}
	})

	t.Run("Timer does not exist", func(t *testing.T) {
		timerMap := make(map[string]*time.Timer)
		key := "non_existent_timer"

		// Call the function with a non-existent timer key
		stopTimerIfExists(timerMap, key)

		// No panic or error should occur, and the map should remain empty
		if len(timerMap) != 0 {
			t.Errorf("Expected timer map to be empty, but got %d elements", len(timerMap))
		}
	})
}
