package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	Laps        int
	LapLen      float64
	PenaltyLen  float64
	FiringLines int
	Start       string
	StartDelta  string
}

func checkErr(err error) error {
	if err != nil {
		return err
	}
	return err
}

// читаем и парсим JSON-файл
func LoadConf(input string) (*Config, error) {
	data, err := os.ReadFile(input)
	checkErr(err)

	var conf Config
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func (c *Config) StartDeltaDuration() (time.Duration, error) {
	t, err := time.Parse("15:04:05", c.StartDelta)
	if err != nil {
		return 0, err
	}
	return time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute + time.Duration(t.Second())*time.Second, nil
}
