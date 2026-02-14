package main

import (
	"github.com/myacey/selectel-logcheck/pkg/logcheck"
	"golang.org/x/tools/go/analysis"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		logcheck.Analyzer,
	}, nil
}
