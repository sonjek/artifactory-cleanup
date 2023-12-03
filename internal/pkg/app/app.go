package app

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/sonjek/artifactory-cleanup/internal/app/actions"
	"github.com/sonjek/artifactory-cleanup/internal/app/params"
)

type App struct {
	Params  *params.QueryParams
	Type    string
	DryRun  bool
	Cleaner actions.ArtifactoryCleaner
	Items   []string
}

func init() {
	log.SetOutput(os.Stdout)
}

func New() (*App, error) {
	// Init CLI params
	cliParams, err := params.InitCLIParams()
	if err != nil {
		return nil, errors.New("parce params: " + err.Error())
	}

	// Init ENV vars
	envVars, err := params.GetEnvVars()
	if err != nil {
		return nil, errors.New("env var evaluation: " + err.Error())
	}

	// Parce config file
	configFile, err := params.ParceConfigFile(cliParams.ConfigFile)
	if err != nil {
		msg := fmt.Sprintf("parce config file %s: %v", cliParams.ConfigFile, err)
		return nil, errors.New(msg)
	}

	endpointRST := envVars.Domain + "/artifactory"
	endpointAQL := endpointRST + "/api/search/aql"

	a := &App{}
	a.Type = configFile.Type
	a.DryRun = cliParams.DryRun

	a.Params = &params.QueryParams{
		ArtifactoryRSTEndpoint: endpointRST,
		ArtifactoryAQLEndpoint: endpointAQL,
		Authorization:          envVars.Authorization,
		RepoNames:              configFile.Repos,
		PatternsArray:          configFile.CleanupPatterns,
		ExcludeArray:           configFile.ExcludePatterns,
		FilterRules:            configFile.Rules,
		DeleteBefore:           cliParams.DeleteBefore,
		Limit:                  cliParams.Limit,
		ShowAQL:                cliParams.ShowAQL,
	}

	switch a.Type {
	case "docker":
		a.Cleaner = actions.Docker{
			Params:        a.Params,
			FilterByField: "path",
			DryRun:        a.DryRun,
		}
	case "file":
		a.Cleaner = actions.File{
			Params:        a.Params,
			FilterByField: "name",
			DryRun:        a.DryRun,
		}
	default:
		return nil, fmt.Errorf("unknown cleaner type: \033[33m%s\033[0m", a.Type)
	}

	return a, nil
}

func (a *App) CollectItems() error {
	items, err := a.Cleaner.Collect()
	if err != nil {
		return err
	}

	a.Items = items

	return nil
}

func (a *App) Clean() error {
	a.Cleaner.Remove(a.Items)
	return nil
}
