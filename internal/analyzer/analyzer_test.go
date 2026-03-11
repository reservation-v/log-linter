package analyzer

import (
	"testing"

	"github.com/reservation-v/log-linter/internal/config"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "a")
}

func TestAnalyzerDisableEnglish(t *testing.T) {
	testdata := analysistest.TestData()

	cfg := config.Default()
	cfg.English = false

	analysistest.Run(t, testdata, New(cfg), "disabledenglish")
}

func TestAnalyzerSuggestedFixes(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "fixes")
}
