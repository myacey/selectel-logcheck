package main

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/myacey/selectel-logcheck/pkg/logcheck"
	"golang.org/x/tools/go/analysis"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	var cfg logcheck.Config
	if conf != nil {
		if err := mapstructure.Decode(conf, &cfg); err != nil {
			return nil, err
		}
	}

	logcheck.ApplyConfig(cfg)

	return []*analysis.Analyzer{
		logcheck.Analyzer,
	}, nil
}
