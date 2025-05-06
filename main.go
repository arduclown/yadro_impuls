package main

import (
	"fmt"
	"sort"

	"github.com/arduclown/yadro_impuls/competitor"
	"github.com/arduclown/yadro_impuls/config"
	"github.com/arduclown/yadro_impuls/event"
	"github.com/arduclown/yadro_impuls/report"
)

func main() {
	conf, err := config.LoadConf("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	events, err := event.LoadEvent("events")
	if err != nil {
		fmt.Printf("Error loading events: %v\n", err)
		return
	}

	competitors := make(map[int]*competitor.Competitor)
	sort.Slice(events, func(i, j int) bool {
		return events[i].Time.Before(events[j].Time)
	})

	logEvents := events
	for _, e := range events {
		if _, exists := competitors[e.CompetitorID]; !exists {
			competitors[e.CompetitorID] = competitor.NewCompetitor(e.CompetitorID)
		}
		comp := competitors[e.CompetitorID]
		outgoing := comp.FightForVictoryCheck(e, conf)
		logEvents = append(logEvents, outgoing...)
	}

	sort.Slice(logEvents, func(i, j int) bool {
		if logEvents[i].Time.Equal(logEvents[j].Time) {
			return logEvents[i].EventID < logEvents[j].EventID
		}
		return logEvents[i].Time.Before(logEvents[j].Time)
	})

	err = event.SaveEventLog(logEvents, "output.log")
	if err != nil {
		fmt.Printf("Error saving event log: %v\n", err)
		return
	}

	var compList []*competitor.Competitor
	for _, c := range competitors {
		compList = append(compList, c)
	}
	sort.Slice(compList, func(i, j int) bool {
		if compList[i].NotFinished || !compList[i].Started {
			return false
		}
		if compList[j].NotFinished || !compList[j].Started {
			return true
		}
		return compList[i].FinishReal.Sub(compList[i].StartReal) < compList[j].FinishReal.Sub(compList[j].StartReal)
	})

	fmt.Println("\nResulting table:")
	fmt.Println(report.GenerateReport(compList))
}
