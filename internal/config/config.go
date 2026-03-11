package config

type Config struct {
	Lowercase              bool
	English                bool
	Symbols                bool
	Sensitive              bool
	ExtraSensitiveKeywords []string
}

func Default() Config {
	return Config{
		Lowercase: true,
		English:   true,
		Symbols:   true,
		Sensitive: true,
	}
}
