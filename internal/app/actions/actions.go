package actions

import (
	"fmt"
	"log"
	"strings"

	"github.com/sonjek/artifactory-cleanup/internal/app/params"
	"github.com/sonjek/artifactory-cleanup/internal/app/rest"
)

type ArtifactoryCleaner interface {
	Collect() ([]string, error)
	Remove(paths []string)
}

func doQueryAndParce(field string, params *params.QueryParams) ([]string, error) {
	if params.ShowAQL {
		log.Printf("AQL query: \n%s", params.AQLQuery)
	}

	values, err := rest.DoAQLQuery(params.ArtifactoryAQLEndpoint, params.AQLQuery, params.Authorization)
	if err != nil {
		return nil, fmt.Errorf("error on get data: %v", err)
	}

	paths := []string{}
	if len(values) > 0 {
		for _, value := range values {
			path := value.Get(field).String()
			repo := value.Get("repo").String()
			paths = append(paths, fmt.Sprintf("%s/%s", repo, path))
		}
	}

	repo := "all repositories"
	if len(params.RepoNames) > 0 {
		repo = fmt.Sprintf("repo: \033[33m%s\033[0m", strings.Join(params.RepoNames, ", "))
	}

	log.Printf("Got \033[33m%d\033[0m items older than \033[33m%s\033[0m for removal from %s\n",
		len(paths), strings.ReplaceAll(params.DeleteBefore, "w", " week"), repo)
	return paths, nil
}

func clean(params *params.QueryParams, paths []string, dryRun bool) {
	if len(paths) < 1 {
		log.Println("No data for removal")
		return
	}

	log.Println("Start removal:")

	totalElements := len(paths)
	for i, path := range paths {
		statusBar := fmt.Sprintf("\033[33m[%d/%d]\033[0m", (i + 1), totalElements)
		finalPath := fmt.Sprintf("%s/%s", params.ArtifactoryRSTEndpoint, path)

		if dryRun {
			log.Printf("%s [DryRun] \033[32m%s\033[0m\n", statusBar, finalPath)
			continue
		}

		log.Printf("%s [Delete] \033[32m%s\033[0m\n", statusBar, finalPath)
		if _, err := rest.DoRSTQuery("DELETE", finalPath, "", params.Authorization); err != nil {
			log.Printf("Error on clean data: %v", err)
		}
	}
}
