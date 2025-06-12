package main

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/jochem11/inventory-system-back/education"
	"github.com/jochem11/inventory-system-back/graphql/generated"
)

type Server struct {
	educationClient *education.Client
}

func (s *Server) Mutation() generated.MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() generated.QueryResolver {
	return &queryResolver{
		server: s,
	}
}

func (s *Server) Subscription() generated.SubscriptionResolver {
	return &subscriptionResolver{
		server: s,
	}
}

func NewGraphQLServer(educationURL string) (*Server, error) {
	educationClient, err := education.NewClient(educationURL)
	if err != nil {
		return nil, err
	}

	return &Server{
		educationClient,
	}, nil
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: s,
	})
}
