package main

import (
	"github.com/myacey/selectel-logcheck/pkg/logcheck"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(logcheck.Analyzer)
}
