package params

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepare() {
	os.Unsetenv("ARTIFACTORY_DOMAIN")
	os.Unsetenv("ARTIFACTORY_USERNAME")
	os.Unsetenv("ARTIFACTORY_PASSWORD")
}

func TestMain(m *testing.M) {
	prepare()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGetEnvVarsMissedEnvs(t *testing.T) {
	_, err := GetEnvVars()
	expected := "you must set: \033[33mARTIFACTORY_DOMAIN, ARTIFACTORY_USERNAME, ARTIFACTORY_PASSWORD\033[30m"
	assert.Error(t, err, expected)
}

func TestGetEnvVarsMissedEnv(t *testing.T) {
	err := os.Setenv("ARTIFACTORY_DOMAIN", "test-domain")
	assert.NoError(t, err)

	err = os.Setenv("ARTIFACTORY_USERNAME", "test-username")
	assert.NoError(t, err)

	_, err = GetEnvVars()
	expected := "you must set: \033[33mARTIFACTORY_PASSWORD\033[30m"
	assert.Error(t, err, expected)
}

func TestGetEnvVarsDefinedValidDomain(t *testing.T) {
	err := os.Setenv("ARTIFACTORY_DOMAIN", "test-domain")
	assert.NoError(t, err)

	err = os.Setenv("ARTIFACTORY_USERNAME", "test-username")
	assert.NoError(t, err)

	err = os.Setenv("ARTIFACTORY_PASSWORD", "test-password")
	assert.NoError(t, err)

	envVars, err := GetEnvVars()
	assert.NoError(t, err)

	expectedDomain := "test-domain"
	assert.Equal(t, expectedDomain, envVars.Domain)
}

func TestGetEnvVarsDefinedValidAuthorization(t *testing.T) {
	err := os.Setenv("ARTIFACTORY_DOMAIN", "test-domain")
	assert.NoError(t, err)

	err = os.Setenv("ARTIFACTORY_USERNAME", "test-username")
	assert.NoError(t, err)

	err = os.Setenv("ARTIFACTORY_PASSWORD", "test-password")
	assert.NoError(t, err)

	envVars, err := GetEnvVars()
	assert.NoError(t, err)

	expectedAuth := "Basic dGVzdC11c2VybmFtZTp0ZXN0LXBhc3N3b3Jk"
	assert.Equal(t, expectedAuth, envVars.Authorization)
}

func TestGetEnvVarValueOrFailDefinedEnv(t *testing.T) {
	key := "MY_ENV_VAR"
	expectedValue := "some_value"
	os.Setenv(key, expectedValue)

	result, err := getEnvVarValueOrFail(key)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, result)
}

func TestGetEnvVarValueOrFailMissedEnv(t *testing.T) {
	key := "NON_EXISTENT_ENV_VAR"
	value, err := getEnvVarValueOrFail(key)
	assert.Error(t, err, key)
	assert.Empty(t, value)
}
