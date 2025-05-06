package competitor

import (
	"time"

	"github.com/arduclown/yadro_impuls/config"
	"github.com/arduclown/yadro_impuls/event"
)

type Competitor struct {
	ID         int
	Registered bool
	StartDraw  time.Time // назначенное время старта при рег. ID = 2
	StartReal  time.Time // реальное время старта ID = 4
	Started    bool      // участник стартовал ID = 4

	FiringRange int // стрельбище (его номер)

	Hits  int // кол-во поподаний
	Shots int // кол-во выстрелов

	PenaltyL     int           // кол-во штрафных кругов
	PenaltyTime  time.Duration // время на шрафном крге
	SpeedPenatly float64       // скорость на шрафном

	CurLap   int // номер текущего круга
	SpeedLap float64

	Finished   bool
	FinishReal time.Time

	CurTime time.Time
}

func (c *Competitor) FightForVictoryCheck(ev event.Event, conf config.Config) {
	c.CurTime = ev.Time

	switch ev.EventID {
	case 1: // регистрация
		c.Registered = true
	case 2: // время назначенное жеребьёвкой
		sTime, _ := time.Parse("15:04:05.000", ev.ExtraParams[0])
		c.StartDraw = sTime
	case 3:
	case 4: // реальный старт
		c.Started = true
		c.StartReal = ev.Time
	case 5: // участник на стрельбище
		c.Hits = 0         // обнулили счетчик на n-ом стрельбище
		c.FiringRange += 1 // кол-во посещенных зон увличили
	case 6: // участник попал в цель
		c.Hits += 1
		c.Shots += 1
	case 7: // участник покинул стрельбище
		// тут считаем сколько штрафных кругов ему накинут
		misses := 5 - c.Hits%5 // пропущенные мишени
		if misses != 0 {
			c.PenaltyL += misses
			c.Shots += misses
		}

	case 8:

	case 9: // участник покинул штрафной круг
		duration := ev.Time.Sub(c.CurTime)
		if c.PenaltyL > 0 {
			c.PenaltyTime = duration

			// скорость для всех штрафных кругов
			c.SpeedPenatly = (conf.PenaltyLen * float64(c.PenaltyL)) / duration.Seconds()
		}

	case 10:
		c.CurLap += 1
		//duration := ev.Time.Sub(c.StartReal)
		if c.CurLap == conf.Laps {
			c.Finished = true
			c.FinishReal = c.CurTime
		}

	}

}
