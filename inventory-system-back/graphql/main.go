//go:generate go run github.com/99designs/gqlgen

package main

import (
	"github.com/99designs/gqlgen/handler"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	EducationUDL string `envconfig:"EDUCATION_SERVICE_URL"`
}

func main() {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	s, err := NewGraphQLServer(cfg.EducationUDL)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/graphql", handler.GraphQL(s.ToExecutableSchema()))
	http.Handle("/playground", handler.Playground("jochem11", "/graphql"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
