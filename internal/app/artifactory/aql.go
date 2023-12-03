package artifactory

// Generate query like:
//
//		items.find({
//		    "repo": "REPO_NAME",
//			"name": { "$match": "manifest.json" },
//	     	"$and": [
//				{ "created": { "$before": "1mo" } },
//				{ "modified": { "$before": "1mo" } }
//			],
//		    "$or": [
//		        { "stat.downloads": { "$eq": null } },
//		        { "stat.downloads": { "$eq": 0 } }
//		    ],
//	   	 	"$or": [
//	      	  { "path": { "$match": "*-gc-*" } },
//	     	   { "path": { "$match": "*develop*" } },
//	     	   { "path": { "$match": "*/custom-*" } }
//	    	]
//		}).include("repo", "path", "name").sort({"$desc":["repo","path"]}).limit(10)
//
// $before can be - 1y, 1mo, 1w, 1d

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sonjek/artifactory-cleanup/internal/app/params"
)

func PrepareSearchDockerImagesAQLQuery(filterByField string, params *params.QueryParams) error {
	repo, err := prepareFilterByRepo(params.RepoNames)
	if err != nil {
		return err
	}

	limitByDate, err := prepareFilterByDate(params.DeleteBefore)
	if err != nil {
		return err
	}

	patterns := preparePatterns(filterByField, params)

	var conditions []string
	for _, e := range []string{
		repo,
		`"name": {"$match": "manifest.json"}`,
		limitByDate,
		patterns,
	} {
		if e != "" {
			conditions = append(conditions, e)
		}
	}

	params.AQLQuery = fmt.Sprintf(`items.find({%s}).include("repo", "%s").sort({"$asc":["path"]}).limit(%d)`,
		strings.Join(conditions, ", "), filterByField, params.Limit)
	return nil
}

func PrepareSearchFilesAQLQuery(filterByField string, params *params.QueryParams) error {
	repo, err := prepareFilterByRepo(params.RepoNames)
	if err != nil {
		return err
	}

	limitByDate, err := prepareFilterByDate(params.DeleteBefore)
	if err != nil {
		return err
	}

	patterns := preparePatterns(filterByField, params)

	var conditions []string
	for _, e := range []string{repo, limitByDate, patterns} {
		if e != "" {
			conditions = append(conditions, e)
		}
	}

	params.AQLQuery = fmt.Sprintf(`items.find({%s}).include("repo", "%s").sort({"$asc":["repo", "name"]}).limit(%d)`,
		strings.Join(conditions, ", "), filterByField, params.Limit)
	return nil
}

func prepareFilterByRepo(repos []string) (string, error) {
	if len(repos) == 0 {
		return "", nil
	}

	return prepareQueryElements([]string{"repo"}, repos, "$eq", "$or"), nil
}

// "$and": [ {"created": {"$before": "2w"}}, {"modified": {"$before": "2w"}} ],
//
// $before can be - 1y, 1mo, 1w, 1d
func prepareFilterByDate(deleteBefore string) (string, error) {
	// return empty string if deleteBefore is empty
	if deleteBefore == "" {
		return "", nil
	}

	regex := regexp.MustCompile(`^\d+(y|mo|w|d)$`)
	if !regex.MatchString(deleteBefore) {
		return "", fmt.Errorf("invalid date value: %s. Should be one of 1y, 1mo, 1w, 1d", deleteBefore)
	}

	return prepareQueryElements([]string{"created", "modified"}, []string{deleteBefore}, "$before", "$and"), nil
}

func preparePatterns(filed string, params *params.QueryParams) string {
	var patterns string

	if len(params.PatternsArray) != 0 {
		patterns = prepareQueryElements([]string{filed}, params.PatternsArray, "$match", "$or")
	}

	if len(params.ExcludeArray) != 0 {
		if patterns != "" {
			patterns += ", "
		}
		patterns += prepareQueryElements([]string{filed}, params.ExcludeArray, "$nmatch", "$and")
	}

	return patterns
}

func prepareQueryElements(keys []string, values []string, elcondition string, grcondition string) string {
	// return empty string if keys or values is empty
	if len(keys) == 0 || len(values) == 0 {
		return ""
	}

	// build array of conditions like
	// {"created": {"$before": "1w"}} or {"path": {"$match": "*develop*"}}
	var queryParts []string
	for _, key := range keys {
		for _, value := range values {
			queryParts = append(queryParts,
				fmt.Sprintf(`{"%s": {"%s": "%s"}}`, key, elcondition, value))
		}
	}

	// join queryParts and wrap to group condition like
	// "$and": [ {"created": {"$before": "1w"}}, {"modified": {"$before": "1w"}} ]
	return fmt.Sprintf(`"%s": [ %s ]`, grcondition, strings.Join(queryParts, ", "))
}
