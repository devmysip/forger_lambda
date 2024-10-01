package utilis

import (
	"fmt"
	"time"
)

func GetCurrentTime() string {
	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return time.Now().Format(time.RFC3339)
	}
	currentTime := time.Now().In(location)
	userUpdatedAt := currentTime.Format(time.RFC3339)

	return userUpdatedAt

}

func GetLocalTime() time.Time {
	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return time.Now()
	}

	now := time.Now().In(istLocation)
	return now
}

func DaysSinceDate(updatedAt string) (int, error) {
	const layout = "2006-01-02T15:04:05Z07:00"
	updatedTime, err := time.Parse(layout, updatedAt)
	if err != nil {
		return 0, fmt.Errorf("failed to parse `updated_at`: %v", err)
	}

	today := time.Now()
	duration := today.Sub(updatedTime)

	days := int(duration.Hours() / 24)

	return days, nil
}
