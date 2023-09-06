package calendar

import (
	"fmt"
	ics "github.com/arran4/golang-ical"
	"net/http"
	"regexp"
	"strings"
)

func GetCalendar(client *http.Client, re *regexp.Regexp, calendarUrl string) ([]*ics.Calendar, error) {
	r, err := client.Get(calendarUrl)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if r.StatusCode == 200 {
		cal, err := ics.ParseCalendar(r.Body)
		if err != nil {
			return nil, err
		}

		calendars := make([]*ics.Calendar, 4)
		calendars[0] = filterCalendar(cal, re, "1", "1")
		calendars[1] = filterCalendar(cal, re, "1", "2")
		calendars[2] = filterCalendar(cal, re, "2", "3")
		calendars[3] = filterCalendar(cal, re, "2", "4")

		return calendars, nil
	} else {
		return nil, fmt.Errorf("invalid status code %d", r.StatusCode)
	}
}

func filterCalendar(cal *ics.Calendar, re *regexp.Regexp, tp string, td string) *ics.Calendar {
	newCal := ics.NewCalendar()

	for _, e := range cal.Events() {
		description := strings.Replace(e.GetProperty(ics.ComponentPropertyDescription).Value, "\\n", "\n", -1)
		matches := re.FindAllStringSubmatch(description, -1)
		for _, match := range matches {
			if match[3] == "ALT" && (match[1] == "" && match[2] == "" || match[1] == "TP" && match[2] == tp || match[1] == "TD" && match[2] == td) {
				newCal.AddVEvent(e)
				break
			}
		}
	}

	return newCal
}
