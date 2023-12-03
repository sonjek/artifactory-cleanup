package actions

import (
	"github.com/sonjek/artifactory-cleanup/internal/app/artifactory"
	"github.com/sonjek/artifactory-cleanup/internal/app/params"
)

type File struct {
	Params        *params.QueryParams
	FilterByField string
	DryRun        bool
}

func (f File) Collect() ([]string, error) {
	err := artifactory.PrepareSearchFilesAQLQuery(f.FilterByField, f.Params)
	if err != nil {
		return nil, err
	}

	return doQueryAndParce(f.FilterByField, f.Params)
}

func (f File) Remove(paths []string) {
	clean(f.Params, paths, f.DryRun)
}
