package format

import "fmt"

type Seconds int

func (s Seconds) ToMinutes() int {
	return int(s) / 60
}

func (s Seconds) ToHours() float64 {
	return float64(s) / 3600
}

func (s Seconds) ToTimeString() string {
	hours := int(s) / 3600
	minutes := (int(s) % 3600) / 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

