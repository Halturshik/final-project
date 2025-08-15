package api

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dStart string, repeat string) (string, error) {
	if repeat == "" {
		return "", nil
	}

	startDate, err := time.Parse(DateFormat, dStart)
	if err != nil {
		return "", fmt.Errorf("не удалось разобрать datestart: %w", err)
	}

	parts := strings.Fields(repeat)

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("неверный формат правила days")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("интервал должен быть от 1 до 400")
		}

		now = now.Truncate(24 * time.Hour)
		startDate = startDate.Truncate(24 * time.Hour)

		date := startDate
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	case "y":
		if len(parts) != 1 {
			return "", errors.New("неверный формат правила year")
		}

		day := startDate.Day()
		month := startDate.Month()
		year := startDate.Year()

		for {
			year++
			date := time.Date(year, month, day, 0, 0, 0, 0, startDate.Location())
			if date.Month() != month {
				date = time.Date(year, time.March, 1, 0, 0, 0, 0, startDate.Location())
			}
			if afterNow(date, now) {
				return date.Format(DateFormat), nil
			}
		}

	case "w":
		if len(parts) != 2 {
			return "", errors.New("некорректное правило weeks")
		}

		daysStr := strings.Split(parts[1], ",")
		validWeekdays := make(map[time.Weekday]bool)
		for _, ds := range daysStr {
			ds = strings.TrimSpace(ds)
			n, err := strconv.Atoi(ds)
			if err != nil || n < 1 || n > 7 {
				return "", errors.New("некорректный день недели")
			}
			var wd time.Weekday
			if n == 7 {
				wd = time.Sunday
			} else {
				wd = time.Weekday(n)
			}
			validWeekdays[wd] = true
		}

		for i := 1; i <= 7; i++ {
			candidate := now.AddDate(0, 0, i)
			if validWeekdays[candidate.Weekday()] && afterNow(candidate, now) {
				return candidate.Format(DateFormat), nil
			}
		}
		return "", errors.New("не удалось найти подходящую дату по правилу weeks")

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("некорректный формат правила months")
		}

		daysStr := strings.Split(parts[1], ",")
		validDays := make(map[int]bool)
		for _, ds := range daysStr {
			ds = strings.TrimSpace(ds)
			n, err := strconv.Atoi(ds)
			if err != nil || n == 0 || n < -2 || n > 31 {
				return "", errors.New("некорректный день месяца")
			}
			validDays[n] = true
		}

		validMonths := make(map[int]bool)
		if len(parts) == 3 {
			monthsStr := strings.Split(parts[2], ",")
			for _, ms := range monthsStr {
				ms = strings.TrimSpace(ms)
				mn, err := strconv.Atoi(ms)
				if err != nil || mn < 1 || mn > 12 {
					return "", errors.New("некорректный месяц")
				}
				validMonths[mn] = true
			}
		}

		now = now.Truncate(24 * time.Hour)
		startDate = startDate.Truncate(24 * time.Hour)

		base := now
		if startDate.After(now) {
			base = startDate
		}
		today := base.Day()

		var dayKeys []int
		for k := range validDays {
			dayKeys = append(dayKeys, k)
		}
		sort.Ints(dayKeys)

		const maxMonthsCheck = 48

		var bestDate *time.Time

		for monthOffset := 0; monthOffset <= maxMonthsCheck; monthOffset++ {
			candidateMonth := base.AddDate(0, monthOffset, 0)
			year, month := candidateMonth.Year(), candidateMonth.Month()

			if len(validMonths) > 0 && !validMonths[int(month)] {
				continue
			}

			lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, candidateMonth.Location()).Day()
			secondLastDay := lastDay - 1

			for _, d := range dayKeys {
				var day int
				if d > 0 {
					day = d
				} else if d == -1 {
					day = lastDay
				} else if d == -2 {
					day = secondLastDay
				} else {
					continue
				}

				if day > lastDay {
					continue
				}

				if monthOffset == 0 && day <= today {
					continue
				}

				candidateDate := time.Date(year, month, day, 0, 0, 0, 0, candidateMonth.Location())

				if afterNow(candidateDate, now) {
					if bestDate == nil || candidateDate.Before(*bestDate) {
						bestDate = &candidateDate
					}
				}
			}
		}

		if bestDate != nil {
			return bestDate.Format(DateFormat), nil
		}

		return "", errors.New("не удалось найти подходящую дату по правилу months")

	default:
		return "", errors.New("неподдерживаемый формат правила")
	}

}
