package tui

import (
	// _ "github.com/goccy/go-graphviz"
	"github.com/sheb-gregor/angiograph/analysis"
)

func Run(path string, patterns []string) error {
	a := analysis.AnalyseOpts{
		Deps:      false,
		Private:   false,
		PrintJSON: false,
		Path:      path,
		Patterns:  patterns,
	}

	return a.Run()
}
