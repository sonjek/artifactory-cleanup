package params

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"
)

type Params struct {
	DeleteBefore string
	Limit        int
	DryRun       bool
	ShowAQL      bool
	ConfigFile   string
}

func InitCLIParams() (*Params, error) {
	configParam := flag.String("config", "", "Path to the configuration file. (mandatory)")
	beforeParam := flag.Int("before", 2, "Set the number of weeks before today to start cleaning. Delete artifacts created/modified before this date.")
	limitParam := flag.Int("limit", 20000, "Set the number of artifacts for clean up.")
	showAQLQueryParam := flag.Bool("showAQL", false, "Set this flag to true to show AQL query. (default false)")
	destroyParam := flag.Bool("destroy", false, "Set this flag to true to allow removal (disable dryRun). (default false)")

	flag.Usage = func() {
		fmt.Printf("Usage of Artifactory cleanup tool:\n\n")
		fmt.Println(`Mandatory environement variables:
  - ARTIFACTORY_DOMAIN
	Set Artifactory address like: https://artifactory.example.com (export ARTIFACTORY_DOMAIN=https://artifactory.example.com) for local run
  - ARTIFACTORY_USERNAME
	Set Artifactory username like: artifactory (export ARTIFACTORY_USERNAME=artifactory) for local run
  - ARTIFACTORY_PASSWORD
	Set Artifactory username like: password (export ARTIFACTORY_PASSWORD=password) for local run`)
		fmt.Printf("\nParams:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// configParam
	if *configParam == "" {
		return nil, errors.New("you must specify configuration file using --config=<filename>")
	}

	// destroyParam
	dryRun := !*destroyParam
	var mode string
	if dryRun {
		mode = "\033[32m[DryRun]\033[0m Search"
	} else {
		mode = "\033[31m[Delete]\033[0m Remove"
	}

	// Don't allow to remove data that created earler than 2 week
	minimalSafeBefore := 2
	if !dryRun && *beforeParam < minimalSafeBefore {
		return nil, fmt.Errorf("'before' parameter is %d and must be %d or greater", *beforeParam, minimalSafeBefore)
	}
	deleteBefore := fmt.Sprintf("%dw", *beforeParam)

	log.Printf("%s older then: \033[33m%s\033[0m, limit: \033[33m%d\033[0m\n",
		mode, strings.ReplaceAll(deleteBefore, "w", " week"), *limitParam)

	params := &Params{
		DeleteBefore: deleteBefore,
		Limit:        *limitParam,
		DryRun:       dryRun,
		ShowAQL:      *showAQLQueryParam,
		ConfigFile:   *configParam,
	}

	return params, nil
}
