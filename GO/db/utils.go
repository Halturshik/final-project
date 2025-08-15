package db

import "time"

func parseSearchDate(search string) (string, bool) {
	t, err := time.Parse("02.01.2006", search)
	if err != nil {
		return "", false
	}
	return t.Format(DateFormat), true
}
