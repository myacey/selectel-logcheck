package logcheck

import (
	"log"
	"strings"
)

type Config struct {
	CheckLowercase bool `mapstructure:"check-lowercase"`
	CheckEnglish   bool `mapstructure:"check-english"`
	CheckSpecial   bool `mapstructure:"check-special"`
	CheckSensitive bool `mapstructure:"check-sensitive"`

	LogFuncs       []string
	SensitiveWords []string `mapstructure:"sensitive-words"`
}

var defaultConfig = Config{
	CheckLowercase: true,
	CheckEnglish:   true,
	CheckSpecial:   true,
	CheckSensitive: true,

	LogFuncs: []string{
		"Print", "Println", "Printf",
		"Sprint", "Sprintln", "Sprintf",
		"Info", "Infof",
		"Warn", "Warnf",
		"Error", "Errorf",
		"Debug", "Debugf",
		"Fatal", "Fatalln", "Fatalf",
		"Panic", "Panicf", "Panicln",
	},
	SensitiveWords: []string{
		"password",
		"pass",
		"token",
		"api_key",
		"apikey",
		"api-key",
		"secret",
		"jwt",
		"bearer",
	},
}

func ApplyConfig(cfg Config) {
	config.CheckLowercase = cfg.CheckLowercase
	config.CheckEnglish = cfg.CheckEnglish
	config.CheckSpecial = cfg.CheckSpecial
	config.CheckSensitive = cfg.CheckSensitive

	if len(cfg.SensitiveWords) > 0 {
		config.SensitiveWords = append(config.SensitiveWords, cfg.SensitiveWords...)
	}

	log.Printf("CONFIG: %+v", cfg)
}

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(v string) error {
	for _, part := range strings.Split(v, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		*s = append(*s, part)
	}
	return nil
}
