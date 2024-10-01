package utilis

import "time"

func GetCurrentTime() string {
	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return time.Now().Format(time.RFC3339)
	}
	currentTime := time.Now().In(location)
	userUpdatedAt := currentTime.Format(time.RFC3339)

	return userUpdatedAt

}
