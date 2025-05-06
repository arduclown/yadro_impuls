package event

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type EventID int

type Event struct {
	Time         time.Time
	EventID      EventID
	CompetitorID int
	ExtraParams  []string
}

const (
	Registered   EventID = 1
	StartTime    EventID = 2
	StartLine    EventID = 3
	Started      EventID = 4
	FiringRange  EventID = 5
	HitTarget    EventID = 6
	LeftFR       EventID = 7
	EnterPenalty EventID = 8
	LeftPentaly  EventID = 9
	EndMainLap   EventID = 10
	NotContinue  EventID = 11
	Disqualofies EventID = 32
	Finished     EventID = 33
)

func checkErr(err error) error {
	if err != nil {
		return err
	}
	return err
}

// парсим отдельные состояния
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
	e.EventID = EventID(eID)

	cID, err := strconv.Atoi(residual[1])
	checkErr(err)
	e.CompetitorID = cID

	if len(residual) > 2 {
		e.ExtraParams = residual[2:]
	}

	return e, nil
}

// читаем из файла
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

// формирование лога
func (e Event) FEvent() string {
	switch e.EventID {
	case Registered:
		return fmt.Sprintf("[%s] The competitor(%d) registered", e.Time.Format("15:04:05.000"), e.CompetitorID)
	case StartTime:
		return fmt.Sprintf("[%s] The start time for competitor(%d) was set by a draw to %s", e.Time.Format("15:04:05.000"), e.CompetitorID, e.ExtraParams[0])
	case StartLine:
		return fmt.Sprintf("[%s] The competitor(%d) is on the start line", e.Time.Format("15:04:05.000"), e.CompetitorID)
	case Started:
		return fmt.Sprintf("[%s] The competitor(%d) has started", e.Time.Format("15:04:05.000"), e.CompetitorID)
	case FiringRange:
		return fmt.Sprintf("[%s] The competitor(%d) is on the firing range(%s)", e.Time.Format("15:04:05.000"), e.CompetitorID, e.ExtraParams[0])
	case HitTarget:
		return fmt.Sprintf("[%s] The target(%s) has been hit by competitor(%d)", e.Time.Format("15:04:05.000"), e.ExtraParams[0], e.CompetitorID)
	case LeftFR:
		return fmt.Sprintf("[%s] The competitor(%d) left the firing range", e.Time.Format("15:04:05.000"), e.CompetitorID)
	case EnterPenalty:
		return fmt.Sprintf("[%s] The competitor(%d) entered the penalty laps", e.Time.Format("15:04:05.000"), e.CompetitorID)
	case LeftPentaly:
		return fmt.Sprintf("[%s] The competitor(%d) left the penalty laps", e.Time.Format("15:04:05.000"), e.CompetitorID)
	case EndMainLap:
		return fmt.Sprintf("[%s] The competitor(%d) ended the main lap", e.Time.Format("15:04:05.000"), e.CompetitorID)
	case NotContinue:
		return fmt.Sprintf("[%s] The competitor(%d) can`t continue: %s", e.Time.Format("15:04:05.000"), e.CompetitorID, strings.Join(e.ExtraParams, " "))
	case Disqualofies:
		return fmt.Sprintf("[%s] The competitor(%d) is disqualified", e.Time.Format("15:04:05.000"), e.CompetitorID)
	case Finished:
		return fmt.Sprintf("[%s] The competitor(%d) has finished", e.Time.Format("15:04:05.000"), e.CompetitorID)
	default:
		return fmt.Sprintf("[%s] Unknown event for competitor(%d)", e.Time.Format("15:04:05.000"), e.CompetitorID)
	}
}

func SaveEventLog(events []Event, output string) error {
	file, err := os.Create(output)
	checkErr(err)
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, e := range events {
		_, err := writer.WriteString(e.FEvent() + "\n")
		checkErr(err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error: %v", err)
	}
	return nil
}
