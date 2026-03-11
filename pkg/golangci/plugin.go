package golangci

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/golangci/plugin-module-register/register"
	projectanalyzer "github.com/reservation-v/log-linter/internal/analyzer"
	"github.com/reservation-v/log-linter/internal/config"
	"golang.org/x/tools/go/analysis"
)

const LinterName = "loglinter"

func init() {
	register.Plugin(LinterName, New)
}

var supportedChecks = map[string]func(*config.Config){
	"lowercase": func(cfg *config.Config) {
		cfg.Lowercase = false
	},
	"english": func(cfg *config.Config) {
		cfg.English = false
	},
	"symbols": func(cfg *config.Config) {
		cfg.Symbols = false
	},
	"sensitive": func(cfg *config.Config) {
		cfg.Sensitive = false
	},
}

type Settings struct {
	Disable []string `json:"disable"`
}

type rawSettings struct {
	Disable     []string `json:"disable"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Settings    Settings `json:"settings"`
}

type Plugin struct {
	settings Settings
}

func New(rawSettings any) (register.LinterPlugin, error) {
	settings, err := decodeSettings(rawSettings)
	if err != nil {
		return nil, err
	}

	if err := settings.validate(); err != nil {
		return nil, err
	}

	return &Plugin{settings: settings}, nil
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{projectanalyzer.New(p.settings.config())}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

func (s Settings) config() config.Config {
	cfg := config.Default()

	for _, name := range s.Disable {
		supportedChecks[name](&cfg)
	}

	return cfg
}

func (s Settings) validate() error {
	for _, name := range s.Disable {
		if _, ok := supportedChecks[name]; ok {
			continue
		}

		return fmt.Errorf("unsupported disabled check %q", name)
	}

	return nil
}

func decodeSettings(raw any) (Settings, error) {
	if raw == nil {
		return Settings{}, nil
	}

	settings, err := register.DecodeSettings[Settings](raw)
	if err == nil {
		return settings, nil
	}

	var wrapper rawSettings
	if err := decode(raw, &wrapper); err != nil {
		return Settings{}, err
	}

	if len(wrapper.Settings.Disable) != 0 {
		return wrapper.Settings, nil
	}

	return Settings{Disable: wrapper.Disable}, nil
}

func decode(raw any, target any) error {
	var buffer bytes.Buffer

	if err := json.NewEncoder(&buffer).Encode(raw); err != nil {
		return fmt.Errorf("encoding settings: %w", err)
	}

	decoder := json.NewDecoder(&buffer)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decoding settings: %w", err)
	}

	return nil
}
