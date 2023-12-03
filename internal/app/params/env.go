package params

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
)

type EnvVars struct {
	Authorization string
	Domain        string
}

func GetEnvVars() (*EnvVars, error) {
	missedEnvs := []string{}

	domain, err := getEnvVarValueOrFail("ARTIFACTORY_DOMAIN")
	if err != nil {
		missedEnvs = append(missedEnvs, err.Error())
	}

	username, err := getEnvVarValueOrFail("ARTIFACTORY_USERNAME")
	if err != nil {
		missedEnvs = append(missedEnvs, err.Error())
	}

	password, err := getEnvVarValueOrFail("ARTIFACTORY_PASSWORD")
	if err != nil {
		missedEnvs = append(missedEnvs, err.Error())
	}

	if len(missedEnvs) > 0 {
		msg := fmt.Sprintf("you must set: \033[33m%s\033[30m", strings.Join(missedEnvs, ", "))
		return &EnvVars{}, errors.New(msg)
	}

	// Prepare auth
	basicAuth := username + ":" + password
	creds := base64.StdEncoding.EncodeToString([]byte(basicAuth))
	auth := "Basic " + creds

	return &EnvVars{
		Authorization: auth,
		Domain:        domain,
	}, nil
}

func getEnvVarValueOrFail(envKey string) (string, error) {
	value := os.Getenv(envKey)
	if value == "" {
		return value, errors.New(envKey)
	}
	return value, nil
}
