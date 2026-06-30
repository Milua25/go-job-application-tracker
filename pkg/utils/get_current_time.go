package utils

import "time"

func GetCurrentTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}
