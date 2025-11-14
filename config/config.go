package config

import (
	"encoding/json"
	"os"
)

type DoubleChar struct {
	Single     string `json:"if_single"`
	SingleType string `json:"single_type"`
	Double     string `json:"if_double"`
	DoubleType string `json:"double_type"`
}

type Config struct {
	Lexer struct {
		TokenTypes  []string          `json:"token_types"`
		KeywordMap  map[string]string `json:"keyword_map"`
		SingleChars map[string]string `json:"single_chars"`
		DoubleChars []DoubleChar      `json:"double_chars"`
	} `json:"lexer"`
}

func ConfigFromJson(data []byte) *Config {
	cfg := &Config{}
	err := json.Unmarshal(data, cfg)
	if err != nil {
		return nil
	}

	return cfg
}

func ConfigFromFile(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return ConfigFromJson(data)
}
