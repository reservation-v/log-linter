package golangci

import (
	"github.com/golangci/plugin-module-register/register"
	projectanalyzer "github.com/reservation-v/log-linter/internal/analyzer"
	"golang.org/x/tools/go/analysis"
)

const LinterName = "loglinter"

func init() {
	register.Plugin(LinterName, New)
}

type Settings struct{}

type Plugin struct {
	settings Settings
}

func New(rawSettings any) (register.LinterPlugin, error) {
	settings, err := register.DecodeSettings[Settings](rawSettings)
	if err != nil {
		return nil, err
	}

	return &Plugin{settings: settings}, nil
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{projectanalyzer.Analyzer}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
