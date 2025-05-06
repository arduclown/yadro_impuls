package competitor

import (
	"time"

	"github.com/arduclown/yadro_impuls/config"
	"github.com/arduclown/yadro_impuls/event"
)

type Competitor struct {
	ID         int
	Registered bool

	StartDraw time.Time // назначенное время старта при рег. ID = 2
	StartReal time.Time // реальное время старта ID = 4
	Started   bool      // участник стартовал ID = 4

	FiringRange int // стрельбище (его номер)
	Hits        int // кол-во попаданий
	Shots       int // кол-во выстрелов

	PenaltyL         int           // кол-во штрафных кругов
	PenaltyTime      time.Duration // время на штрафном круге
	PenaltyStartTime time.Time     //начало штрафного
	SpeedPenalty     float64       // скорость на штрафном

	CurLap       int       // номер текущего круга
	Lap          []Laps    // инфа о текущем круге
	LapStartTime time.Time // начало обычного круга

	Finished    bool
	FinishReal  time.Time
	CurTime     time.Time
	NotFinished bool
}

type Laps struct {
	TimeSpent time.Duration
	Speed     float64
}

func NewCompetitor(id int) *Competitor {
	return &Competitor{
		ID:               id,
		Started:          false,
		Finished:         false,
		NotFinished:      false,
		LapStartTime:     time.Time{},
		PenaltyStartTime: time.Time{},
	}
}

func (c *Competitor) FightForVictoryCheck(ev event.Event, conf *config.Config) []event.Event {
	var out []event.Event

	c.CurTime = ev.Time
	switch ev.EventID {
	case event.Registered:
		c.Registered = true
	case event.StartTime:
		sTime, _ := time.Parse("15:04:05.000", ev.ExtraParams[0])
		c.StartDraw = sTime
	case event.StartLine:
	case event.Started:
		c.Started = true
		c.StartReal = ev.Time
		c.LapStartTime = ev.Time
		startDelta, _ := conf.StartDeltaDuration()
		if c.StartDraw.Add(startDelta).Before(ev.Time) {
			c.NotFinished = true
			out = append(out, event.Event{
				Time:         ev.Time,
				EventID:      event.Disqualofies,
				CompetitorID: c.ID,
			})
		}
	case event.FiringRange:
		c.Hits = 0
		c.Shots = 0
		c.FiringRange += 1
	case event.HitTarget:
		if c.Hits < 5 {
			c.Hits += 1
			c.Shots += 1
		}
	case event.LeftFR:
		misses := 5 - c.Hits
		if misses > 0 {
			c.PenaltyL += misses
		}
		c.Shots = 5
	case event.EnterPenalty:
		if c.PenaltyStartTime.IsZero() {
			c.PenaltyStartTime = ev.Time
		}
	case event.LeftPentaly:
		if c.PenaltyL > 0 && !c.PenaltyStartTime.IsZero() {
			duration := ev.Time.Sub(c.PenaltyStartTime)
			if duration > 0 {
				c.PenaltyTime = duration
				c.SpeedPenalty = (conf.PenaltyLen * float64(c.PenaltyL)) / duration.Seconds()
				c.PenaltyL = 0
				c.PenaltyStartTime = time.Time{}
			}
		}
	case event.EndMainLap:
		if !c.LapStartTime.IsZero() {
			c.CurLap += 1
			duration := ev.Time.Sub(c.LapStartTime)
			if duration > 0 {
				c.Lap = append(c.Lap, Laps{
					TimeSpent: duration,
					Speed:     conf.LapLen / duration.Seconds(),
				})
				c.LapStartTime = ev.Time
			}
			if c.CurLap == conf.Laps && !c.NotFinished {
				c.Finished = true
				c.FinishReal = c.CurTime
				out = append(out, event.Event{
					Time:         ev.Time,
					EventID:      event.Finished,
					CompetitorID: c.ID,
				})
			}
		}
	case event.NotContinue:
		c.NotFinished = true
	}

	return out
}
