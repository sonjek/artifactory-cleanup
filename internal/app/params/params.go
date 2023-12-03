package params

type QueryParams struct {
	ArtifactoryRSTEndpoint string
	ArtifactoryAQLEndpoint string
	Authorization          string
	RepoNames              []string
	DeleteBefore           string
	PatternsArray          []string
	ExcludeArray           []string
	FilterRules            []Rule
	Limit                  int
	AQLQuery               string
	ShowAQL                bool
}
