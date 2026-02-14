package logcheck

import "strings"

type Config struct {
	CheckLowercase bool
	CheckEnglish   bool
	CheckSpecial   bool
	CheckSensitive bool

	LogFuncs       []string
	SensitiveWords []string
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

func loadConfig() {
	Analyzer.Flags.BoolVar(&config.CheckLowercase, "lowercase", true, "check lowercase")
	Analyzer.Flags.BoolVar(&config.CheckEnglish, "english", true, "check english")
	Analyzer.Flags.BoolVar(&config.CheckSpecial, "special", true, "check special chars")
	Analyzer.Flags.BoolVar(&config.CheckSensitive, "sensitive", true, "check sensitive data")

	additionalSensitiveWords := []string{}
	Analyzer.Flags.Var(
		(*stringSlice)(&additionalSensitiveWords),
		"sensitive-words",
		"additional sensitive words",
	)
	config.SensitiveWords = append(config.SensitiveWords, additionalSensitiveWords...)
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
