package internal

import (
	"fmt"
	"math"
	"time"
)

func GetAge(currentTime time.Time, reference time.Time) string {
	diff := int(math.Round(currentTime.Sub(reference).Minutes()))

	if diff < 0 {
		return "-"
	} else if diff == 1 {
		return fmt.Sprintf("%d minute", diff)
	} else if diff < 60 {
		return fmt.Sprintf("%d minutes", diff)
	} else {
		h := diff / 60

		if h == 1 {
			return fmt.Sprintf("%d hour", h)
		}

		return fmt.Sprintf("%d hours", h)
	}
}
