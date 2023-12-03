package actions

import (
	"fmt"

	"github.com/sonjek/artifactory-cleanup/internal/app/artifactory"
	"github.com/sonjek/artifactory-cleanup/internal/app/params"
)

type Docker struct {
	Params        *params.QueryParams
	FilterByField string
	DryRun        bool
}

func (d Docker) Collect() ([]string, error) {
	err := artifactory.PrepareSearchDockerImagesAQLQuery(d.FilterByField, d.Params)
	if err != nil {
		return nil, fmt.Errorf("error on AQL query prepare: %v", err)
	}

	return doQueryAndParce(d.FilterByField, d.Params)
}

func (d Docker) Remove(paths []string) {
	clean(d.Params, paths, d.DryRun)
}
