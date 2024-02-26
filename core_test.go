package bizday_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/fernandosavio/bizday"
)

func TestDateStringValidation(t *testing.T) {
	invalidInputs := []string{
		"",
		"1999-02-29",
		"1999-02-AA",
		"1999-02",
		"1999-13-01",
		"1999-01-32",
		"0001-01-01",
		"24-01-01",
		"2024-1-01",
		"2024-01-1",
	}

	for _, input := range invalidInputs {
		testName := fmt.Sprintf("%#v should be an invalid date", input)
		t.Run(testName, func(t *testing.T) {
			holiday, err := bizday.NewHoliday(input)

			if holiday != nil {
				t.Error("should not return a value, only an error")
			}

			if err == nil {
				t.Error("should raise error")
			}
		})
	}

	validInputs := []string{
		"2024-01-01",
		"2024-02-29",
		"2024-12-31",
	}

	for _, input := range validInputs {
		testName := fmt.Sprintf("%#v should be valid", input)
		t.Run(testName, func(t *testing.T) {
			holiday, err := bizday.NewHoliday(input)

			if err != nil {
				t.Error("should not raise any error")
			}

			if *holiday != bizday.Holiday(input) {
				t.Error("holiday date should be the same as the input")
			}
		})
	}
}

func TestTimeValidation(t *testing.T) {
	t.Run("time.Time zero value should be invalid", func(t *testing.T) {
		value := time.Time{}
		holiday, err := bizday.NewHolidayFromTime(value)

		if holiday != nil {
			t.Error("should not return a value, only an error")
		}

		if err == nil {
			t.Error("should raise error")
		}
	})

	t.Run("Should accept valid values", func(t *testing.T) {
		value, _ := time.Parse(time.DateOnly, "2024-02-29")
		holiday, err := bizday.NewHolidayFromTime(value)

		if err != nil {
			t.Error("should not raise any error")
		}

		if *holiday != "2024-02-29" {
			t.Error("holiday date should be the same as the input")
		}
	})
}

type fakeHolidayService struct {
	holidays map[string]bool
}

func (dh *fakeHolidayService) HasHoliday(date time.Time) bool {
	_, found := dh.holidays[date.Format(time.DateOnly)]
	return found
}

func TestNextBusinessDayNoHolidays(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "After monday, should be tuesday",
			input:    "2024-01-01",
			expected: "2024-01-02",
		},
		{
			name:     "After tuesday, should be wednesday",
			input:    "2024-01-02",
			expected: "2024-01-03",
		},
		{
			name:     "After wednesday, should be thursday",
			input:    "2024-01-03",
			expected: "2024-01-04",
		},
		{
			name:     "After thursday, should be friday",
			input:    "2024-01-04",
			expected: "2024-01-05",
		},
		{
			name:     "After friday, should be monday",
			input:    "2024-01-05",
			expected: "2024-01-08",
		},
		{
			name:     "After saturday, should be monday",
			input:    "2024-01-06",
			expected: "2024-01-08",
		},
		{
			name:     "After sunday, should be monday",
			input:    "2024-01-07",
			expected: "2024-01-08",
		},
	}

	service := fakeHolidayService{}
	registry, err := bizday.NewHolidayRegistry(&service)

	if err != nil {
		t.Error("Failed creating registry")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := time.Parse(time.DateOnly, tt.input)

			if err != nil {
				t.Error("Failed parsin input date")
			}

			got := registry.GetNextBusinessDayFrom(input).Format(time.DateOnly)

			if got != tt.expected {
				t.Errorf("Expected %s but received %s", tt.expected, got)
			}
		})
	}
}

func TestNextBusinessDayWithHolidays(t *testing.T) {
	holidays := map[string]bool{
		"2024-01-01": true,
		"2024-01-03": true,
		"2024-01-05": true,
		"2024-01-08": true,
	}
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "After monday, should be tuesday",
			input:    "2024-01-01",
			expected: "2024-01-02",
		},
		{
			name:     "After tuesday, should be thursday",
			input:    "2024-01-02",
			expected: "2024-01-04",
		},
		{
			name:     "After wednesday, should be thursday",
			input:    "2024-01-03",
			expected: "2024-01-04",
		},
		{
			name:     "After thursday, should be tuesday",
			input:    "2024-01-04",
			expected: "2024-01-09",
		},
		{
			name:     "After friday, should be tuesday",
			input:    "2024-01-05",
			expected: "2024-01-09",
		},
		{
			name:     "After saturday, should be tuesday",
			input:    "2024-01-06",
			expected: "2024-01-09",
		},
		{
			name:     "After sunday, should be tuesday",
			input:    "2024-01-07",
			expected: "2024-01-09",
		},
	}

	service := fakeHolidayService{holidays: holidays}
	registry, err := bizday.NewHolidayRegistry(&service)

	if err != nil {
		t.Error("Failed creating registry")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := time.Parse(time.DateOnly, tt.input)

			if err != nil {
				t.Error("Failed parsin input date")
			}

			got := registry.GetNextBusinessDayFrom(input).Format(time.DateOnly)

			if got != tt.expected {
				t.Errorf("Expected %s but received %s", tt.expected, got)
			}
		})
	}
}

func TestIsBusinessDay(t *testing.T) {
	holidays := map[string]bool{
		"2024-01-01": true,
		"2024-01-03": true,
		"2024-01-05": true,
		"2024-01-08": true,
	}
	tests := []struct {
		input    string
		expected bool
	}{
		{
			input:    "2024-01-01",
			expected: false,
		},
		{
			input:    "2024-01-02",
			expected: true,
		},
		{
			input:    "2024-01-03",
			expected: false,
		},
		{
			input:    "2024-01-04",
			expected: true,
		},
		{
			input:    "2024-01-05",
			expected: false,
		},
		{
			input:    "2024-01-06",
			expected: false,
		},
		{
			input:    "2024-01-07",
			expected: false,
		},
		{
			input:    "2024-01-08",
			expected: false,
		},
		{
			input:    "2024-01-09",
			expected: true,
		},
	}

	service := fakeHolidayService{holidays: holidays}
	registry, err := bizday.NewHolidayRegistry(&service)

	if err != nil {
		t.Error("Failed creating registry")
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%#v should be %v", tt.input, tt.expected)
		t.Run(testName, func(t *testing.T) {
			input, err := time.Parse(time.DateOnly, tt.input)

			if err != nil {
				t.Error("Failed parsin input date")
			}

			got := registry.IsBusinessDay(input)

			if got != tt.expected {
				t.Errorf("Expected %v but received %v", tt.expected, got)
			}
		})
	}
}
