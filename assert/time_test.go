package assert_test

import (
	"testing"
	"time"

	"github.com/neticdk/go-stdlib/assert"
)

func TestTimeAfter(t *testing.T) {
	now := time.Now()
	before := now.Add(-1 * time.Hour)
	after := now.Add(1 * time.Hour)

	tests := []struct {
		name      string
		got       time.Time
		threshold time.Time
		wantPass  bool
	}{
		{"time is after threshold", after, now, true},
		{"time is before threshold", before, now, false},
		{"time is equal to threshold", now, now, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}
			result := assert.TimeAfter(mockT, tt.got, tt.threshold)
			if result != tt.wantPass {
				t.Errorf("TimeAfter() = %v, want %v", result, tt.wantPass)
			}
			if mockT.Failed() != !tt.wantPass {
				t.Errorf("mockT.Failed() = %v, want %v", mockT.Failed(), !tt.wantPass)
			}
		})
	}
}

func TestTimeBefore(t *testing.T) {
	now := time.Now()
	before := now.Add(-1 * time.Hour)
	after := now.Add(1 * time.Hour)

	tests := []struct {
		name      string
		got       time.Time
		threshold time.Time
		wantPass  bool
	}{
		{"time is before threshold", before, now, true},
		{"time is after threshold", after, now, false},
		{"time is equal to threshold", now, now, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}
			result := assert.TimeBefore(mockT, tt.got, tt.threshold)
			if result != tt.wantPass {
				t.Errorf("TimeBefore() = %v, want %v", result, tt.wantPass)
			}
			if mockT.Failed() != !tt.wantPass {
				t.Errorf("mockT.Failed() = %v, want %v", mockT.Failed(), !tt.wantPass)
			}
		})
	}
}

func TestTimeEqual(t *testing.T) {
	now := time.Now()
	other := now.Add(1 * time.Minute)
	sameTimeOtherZone := now.In(time.FixedZone("UTC+2", 2*60*60))

	tests := []struct {
		name     string
		got      time.Time
		want     time.Time
		wantPass bool
	}{
		{"identical times", now, now, true},
		{"different times", now, other, false},
		{"same time different zone", now, sameTimeOtherZone, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}
			result := assert.TimeEqual(mockT, tt.got, tt.want)
			if result != tt.wantPass {
				t.Errorf("TimeEqual() = %v, want %v", result, tt.wantPass)
			}
			if mockT.Failed() != !tt.wantPass {
				t.Errorf("mockT.Failed() = %v, want %v", mockT.Failed(), !tt.wantPass)
			}
		})
	}
}

func TestWithinDuration(t *testing.T) {
	now := time.Now()
	nearTime := now.Add(30 * time.Second)
	farTime := now.Add(5 * time.Minute)

	tests := []struct {
		name     string
		got      time.Time
		want     time.Time
		delta    time.Duration
		wantPass bool
	}{
		{"within duration", nearTime, now, 1 * time.Minute, true},
		{"exactly at boundary", nearTime, now, 30 * time.Second, true},
		{"outside duration", farTime, now, 1 * time.Minute, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}
			result := assert.WithinDuration(mockT, tt.got, tt.want, tt.delta)
			if result != tt.wantPass {
				t.Errorf("WithinDuration() = %v, want %v", result, tt.wantPass)
			}
			if mockT.Failed() != !tt.wantPass {
				t.Errorf("mockT.Failed() = %v, want %v", mockT.Failed(), !tt.wantPass)
			}
		})
	}
}

func TestTimeEqualWithPrecision(t *testing.T) {
	baseTime := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	sameMinute := time.Date(2023, 5, 15, 10, 30, 45, 0, time.UTC)
	differentMinute := time.Date(2023, 5, 15, 10, 31, 0, 0, time.UTC)

	tests := []struct {
		name      string
		got       time.Time
		want      time.Time
		precision time.Duration
		wantPass  bool
	}{
		{"same to the minute", baseTime, sameMinute, time.Minute, true},
		{"different minutes, same hour", baseTime, differentMinute, time.Hour, true},
		{"different minutes", baseTime, differentMinute, time.Minute, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}
			result := assert.TimeEqualWithPrecision(mockT, tt.got, tt.want, tt.precision)
			if result != tt.wantPass {
				t.Errorf("TimeEqualWithPrecision() = %v, want %v", result, tt.wantPass)
			}
			if mockT.Failed() != !tt.wantPass {
				t.Errorf("mockT.Failed() = %v, want %v", mockT.Failed(), !tt.wantPass)
			}
		})
	}
}

func TestWithinTime(t *testing.T) {
	now := time.Now()
	start := now.Add(-1 * time.Hour)
	end := now.Add(1 * time.Hour)
	beforeStart := start.Add(-1 * time.Minute)
	afterEnd := end.Add(1 * time.Minute)

	tests := []struct {
		name     string
		got      time.Time
		start    time.Time
		end      time.Time
		wantPass bool
	}{
		{"within time window", now, start, end, true},
		{"at start boundary", start, start, end, true},
		{"at end boundary", end, start, end, true},
		{"before time window", beforeStart, start, end, false},
		{"after time window", afterEnd, start, end, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}
			result := assert.WithinTime(mockT, tt.got, tt.start, tt.end)
			if result != tt.wantPass {
				t.Errorf("WithinTime() = %v, want %v", result, tt.wantPass)
			}
			if mockT.Failed() != !tt.wantPass {
				t.Errorf("mockT.Failed() = %v, want %v", mockT.Failed(), !tt.wantPass)
			}
		})
	}
}
