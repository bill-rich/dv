package main

import (
	"github.com/bill-rich/dv/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

type dvAnalyzer struct{}

var AnalyzerPlugin dvAnalyzer

func (analyzerStruct *dvAnalyzer) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{analyzer.Analyzer}
}
