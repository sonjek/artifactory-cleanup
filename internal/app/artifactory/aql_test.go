package artifactory

import (
	"testing"

	"github.com/sonjek/artifactory-cleanup/internal/app/params"

	"github.com/stretchr/testify/assert"
)

func TestPrepareQueryElements(t *testing.T) {
	t.Run("EmptyKeys", func(t *testing.T) {
		result := prepareQueryElements([]string{"created"}, []string{}, "$before", "$and")
		assert.Empty(t, result)
	})

	t.Run("EmptyValues", func(t *testing.T) {
		result := prepareQueryElements([]string{}, []string{"1w"}, "$before", "$and")
		assert.Empty(t, result)
	})

	t.Run("TwoKeysOneValue", func(t *testing.T) {
		keys := []string{"created", "modified"}
		values := []string{"1w"}
		elcondition := "$before"
		grcondition := "$and"

		result := prepareQueryElements(keys, values, elcondition, grcondition)
		expected := `"$and": [ {"created": {"$before": "1w"}}, {"modified": {"$before": "1w"}} ]`
		assert.Equal(t, expected, result)
	})
}

func TestPrepareFilterByRepo(t *testing.T) {
	t.Run("EmptyRepos", func(t *testing.T) {
		result, err := prepareFilterByRepo([]string{})
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("NonEmptyRepos", func(t *testing.T) {
		result, err := prepareFilterByRepo([]string{"repo1", "repo2"})
		assert.NoError(t, err)

		expected := `"$or": [ {"repo": {"$eq": "repo1"}}, {"repo": {"$eq": "repo2"}} ]`
		assert.Equal(t, expected, result)
	})
}

func TestPrepareFilterByDate(t *testing.T) {
	t.Run("ValidDeleteBefore1w", func(t *testing.T) {
		result, err := prepareFilterByDate("1w")
		assert.NoError(t, err)

		expected := `"$and": [ {"created": {"$before": "1w"}}, {"modified": {"$before": "1w"}} ]`
		assert.Equal(t, expected, result)
	})

	t.Run("ValidDeleteBefore1mo", func(t *testing.T) {
		result, err := prepareFilterByDate("1mo")
		assert.NoError(t, err)

		expected := `"$and": [ {"created": {"$before": "1mo"}}, {"modified": {"$before": "1mo"}} ]`
		assert.Equal(t, expected, result)
	})

	t.Run("InvalidDeleteBefore", func(t *testing.T) {
		deleteBefore := "1m"
		result, err := prepareFilterByDate(deleteBefore)

		expectedErrorMessage := "invalid date value: " + deleteBefore + ". Should be one of 1y, 1mo, 1w, 1d"
		assert.Error(t, err, expectedErrorMessage)
		assert.Empty(t, result)
	})

	t.Run("EmptyDeleteBefore", func(t *testing.T) {
		result, err := prepareFilterByDate("")
		assert.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestPreparePatterns(t *testing.T) {
	t.Run("PatternsAndExclude", func(t *testing.T) {
		queryParams := &params.QueryParams{
			PatternsArray: []string{"*develop*"},
			ExcludeArray:  []string{"*develop.1"},
		}

		result := preparePatterns("path", queryParams)
		expected := `"$or": [ {"path": {"$match": "*develop*"}} ], "$and": [ {"path": {"$nmatch": "*develop.1"}} ]`
		assert.Equal(t, expected, result)
	})

	t.Run("EmptyPatternsArrayAndExcludeArray", func(t *testing.T) {
		queryParams := &params.QueryParams{
			PatternsArray: []string{},
			ExcludeArray:  []string{},
		}
		result := preparePatterns("path", queryParams)
		assert.Empty(t, result)
	})

	t.Run("EmptyPatternsArray", func(t *testing.T) {
		queryParams := &params.QueryParams{
			PatternsArray: []string{},
			ExcludeArray:  []string{"*test*"},
		}
		result := preparePatterns("path", queryParams)
		expected := `"$and": [ {"path": {"$nmatch": "*test*"}} ]`
		assert.Equal(t, expected, result)
	})

	t.Run("EmptyExcludeArray", func(t *testing.T) {
		queryParams := &params.QueryParams{
			PatternsArray: []string{"*develop*"},
			ExcludeArray:  []string{},
		}

		result := preparePatterns("path", queryParams)
		expected := `"$or": [ {"path": {"$match": "*develop*"}} ]`
		assert.Equal(t, expected, result)
	})
}

func TestPrepareSearchDockerImagesAQLQuery(t *testing.T) {
	params := &params.QueryParams{
		RepoNames:     []string{"repo1"},
		PatternsArray: []string{"*develop*"},
		ExcludeArray:  []string{"*develop.1"},
		DeleteBefore:  "2w",
		Limit:         1,
	}

	err := PrepareSearchDockerImagesAQLQuery("path", params)
	assert.NoError(t, err)

	expectedQuery := `items.find({"$or": [ {"repo": {"$eq": "repo1"}} ], "name": {"$match": "manifest.json"}, "$and": [ {"created": {"$before": "2w"}}, {"modified": {"$before": "2w"}} ], "$or": [ {"path": {"$match": "*develop*"}} ], "$and": [ {"path": {"$nmatch": "*develop.1"}} ]}).include("repo", "path").sort({"$asc":["path"]}).limit(1)`
	assert.Equal(t, expectedQuery, params.AQLQuery)
}

func TestPrepareSearchFilesAQLQuery(t *testing.T) {
	params := &params.QueryParams{
		RepoNames:     []string{"repo1"},
		PatternsArray: []string{"*develop*"},
		ExcludeArray:  []string{"*develop.1"},
		DeleteBefore:  "2w",
		Limit:         1,
	}

	err := PrepareSearchFilesAQLQuery("name", params)
	assert.NoError(t, err)

	expectedQuery := `items.find({"$or": [ {"repo": {"$eq": "repo1"}} ], "$and": [ {"created": {"$before": "2w"}}, {"modified": {"$before": "2w"}} ], "$or": [ {"name": {"$match": "*develop*"}} ], "$and": [ {"name": {"$nmatch": "*develop.1"}} ]}).include("repo", "name").sort({"$asc":["repo", "name"]}).limit(1)`
	assert.Equal(t, expectedQuery, params.AQLQuery)
}
