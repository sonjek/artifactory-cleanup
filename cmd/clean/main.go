package main

import (
	"log"

	"github.com/sonjek/artifactory-cleanup/internal/pkg/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("Error during initialization: %v\n", err)
	}

	if err := a.CollectItems(); err != nil {
		log.Fatalf("Error durring collecting data: %v\n", err)
	}

	if err := a.Clean(); err != nil {
		log.Fatalf("Error durring cleanup of %s: %v\n", a.Params.RepoNames, err)
	}
}
