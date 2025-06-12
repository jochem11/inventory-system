package main

import (
	"context"
	"github.com/jochem11/inventory-system-back/graphql/generated"
)

type subscriptionResolver struct {
	server *Server
}

func (r *subscriptionResolver) LiveCourses(ctx context.Context, pagination *generated.PaginationInput) (<-chan []*generated.Course, error) {
	skip, take := getPaginationBounds(pagination)

	coursesChan, err := r.server.educationClient.LiveCourses(ctx, skip, take)
	if err != nil {
		return nil, err
	}

	ch := make(chan []*generated.Course)

	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case courses, ok := <-coursesChan:
				if !ok {
					return
				}

				gqlCourses := make([]*generated.Course, 0, len(courses))
				for _, c := range courses {
					gqlCourses = append(gqlCourses, &generated.Course{
						ID:        c.ID,
						Name:      c.Name,
						CreatedAt: c.CreatedAt,
						UpdatedAt: c.UpdatedAt,
					})
				}

				select {
				case ch <- gqlCourses:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return ch, nil
}

func (r *subscriptionResolver) LiveClasses(ctx context.Context, pagination *generated.PaginationInput) (<-chan []*generated.Class, error) {
	// Implement similar to LiveCourses
	return nil, nil
}

func getPaginationBounds2(p *generated.PaginationInput) (uint64, uint64) {
	var skip, take uint64 = 0, 50
	if p != nil {
		if p.Skip != nil {
			skip = uint64(*p.Skip)
		}
		if p.Take != nil {
			take = uint64(*p.Take)
		}
	}
	return skip, take
}
