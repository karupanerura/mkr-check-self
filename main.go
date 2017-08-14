package main

import (
	"log"

	"github.com/alecthomas/kingpin"
	"github.com/fatih/color"
	"github.com/mackerelio/mackerel-agent/checks"
	"github.com/mackerelio/mackerel-agent/config"
	colorable "github.com/mattn/go-colorable"
)

var (
	conffile = kingpin.Flag("conf", "Config file path (Configs in this file are over-written by command line options)").Short('c').Default(config.DefaultConfig.Conffile).String()
	verbose  = kingpin.Flag("verbose", "Verbose mode").Short('v').Default("false").Bool()
)

func main() {
	kingpin.Parse()
	logger := log.New(colorable.NewColorableStdout(), "", log.Ldate|log.Ltime|log.LUTC)

	conf, err := config.LoadConfig(*conffile)
	if err != nil {
		logger.Fatal(err)
	}
	// pp.Print(conf)

	for name, config := range conf.CheckPlugins {
		checker := checks.Checker{Name: name, Config: config}
		report := checker.Check()
		logger.Printf("%s: %s\n", statusColordString(report.Status), color.BlueString(checker.String()))
		if *verbose {
			logger.Print(report.Message)
		}
	}
}

func statusColordString(status checks.Status) string {
	switch status {
	case checks.StatusOK:
		return color.GreenString(string(status))
	case checks.StatusWarning:
		return color.YellowString(string(status))
	case checks.StatusCritical:
		return color.RedString(string(status))
	default:
		return color.MagentaString("UNKNWON")
	}
}
