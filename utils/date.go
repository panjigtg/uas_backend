package utils

import (
	"time"
)

func FormatDate(v interface{}) string {
	if v == nil {
		return ""
	}

	switch t := v.(type) {

	case string:
		parsed, err := time.Parse(time.RFC3339, t)
		if err == nil {
			return parsed.Format("2006-01-02")
		}
	
		return t

	case time.Time:
		return t.Format("2006-01-02")
	}

	
	return ""
}


func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
