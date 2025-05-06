package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/arduclown/yadro_impuls/competitor"
)

// функция для преобразования времени в корректное представление по заданию
func FormattingTime(d time.Duration) string {
	hours := int(d / time.Hour)
	d -= time.Duration(hours) * time.Hour
	minutes := int(d / time.Minute)
	d -= time.Duration(minutes) * time.Minute
	seconds := int(d / time.Second)
	milliseconds := int((d % time.Second) / time.Millisecond)

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}

func GenerateReport(comps []*competitor.Competitor) string {
	var result []string
	for _, c := range comps {
		var lapsStr string
		if len(c.Lap) > 0 {
			lapsStr = "["
			for i, lap := range c.Lap {
				if i > 0 {
					lapsStr += ", "
				}
				lapsStr += fmt.Sprintf("{%s, %.3f}", FormattingTime(lap.TimeSpent), lap.Speed)
			}
			lapsStr += "]"
		} else {
			lapsStr = "[{,}]"
		}

		penaltyStr := fmt.Sprintf("{%s, %.3f}", FormattingTime(c.PenaltyTime), c.SpeedPenalty)
		var status string
		if c.NotFinished {
			status = "[NotFinished]"
		} else if !c.Started {
			status = "[NotStarted]"
		} else if c.Finished {
			totalTime := c.FinishReal.Sub(c.StartReal)
			status = FormattingTime(totalTime)
		} else {
			status = "[NotFinished]"
		}

		line := fmt.Sprintf("%s %d %s %s %d/%d", status, c.ID, lapsStr, penaltyStr, c.Hits, c.Shots)
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}
