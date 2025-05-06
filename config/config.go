package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Laps        int
	LapLen      float64
	PenaltyLen  float64
	FiringLines float64
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
