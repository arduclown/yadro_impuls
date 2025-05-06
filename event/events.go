package event

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Time         time.Time
	EventID      int
	CompetitorID int
	ExtraParams  []string
}

func checkErr(err error) error {
	if err != nil {
		return err
	}
	return err
}

func parseEvent(state string) (Event, error) {
	var e Event

	timeEnd := strings.Index(state, "]")
	timePart := state[1:timeEnd]

	t, err := time.Parse("15:04:05.000", timePart)
	checkErr(err)
	e.Time = t

	residual := strings.Fields(state[timeEnd+1:])
	eID, err := strconv.Atoi(residual[0])
	checkErr(err)
	e.EventID = eID

	cID, err := strconv.Atoi(residual[1])
	checkErr(err)
	e.CompetitorID = cID

	if len(residual) > 2 {
		e.ExtraParams = residual[2:]
	}

	return e, nil
}

func LoadEvent(input string) ([]Event, error) {
	var events []Event

	file, err := os.Open(input)
	checkErr(err)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		if l == "" {
			continue
		}

		event, err := parseEvent(l)
		checkErr(err)
		events = append(events, event)
	}

	return events, nil
}
