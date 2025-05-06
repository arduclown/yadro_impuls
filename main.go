package main

import (
	"fmt"

	"github.com/arduclown/yadro_impuls/event"
)

func main() {
	ev, err := event.LoadEvent("events")
	if err != nil {
		fmt.Println("Error!", err)
		return
	}

	for _, e := range ev {
		fmt.Printf("[%s] EventID=%d Competitor=%d ExtraParams=%v\n",
			e.Time.Format("15:04:05.000"), e.EventID, e.CompetitorID, e.ExtraParams)
	}

}
