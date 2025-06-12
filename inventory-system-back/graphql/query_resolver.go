package main

import (
	"context"
	"github.com/jochem11/inventory-system-back/graphql/generated"
	"log"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r queryResolver) Courses(ctx context.Context, pagination *generated.PaginationInput, id *string) ([]*generated.Course, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		r, err := r.server.educationClient.GetCourse(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*generated.Course{{
			ID:        r.ID,
			Name:      r.Name,
			UpdatedAt: r.UpdatedAt,
			CreatedAt: r.CreatedAt,
		}}, nil
	}

	skip, take := getPaginationBounds(pagination)

	coursesList, err := r.server.educationClient.GetCourses(ctx, skip, take)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var courses []*generated.Course
	for _, course := range coursesList {
		courses = append(courses, &generated.Course{
			ID:        course.ID,
			Name:      course.Name,
			UpdatedAt: course.UpdatedAt,
			CreatedAt: course.CreatedAt,
		})
	}
	return courses, nil
}

func (r queryResolver) Classes(ctx context.Context, pagination *generated.PaginationInput, id *string) ([]*generated.Class, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		r, err := r.server.educationClient.GetClass(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return []*generated.Class{{
			ID:       r.ID,
			Name:     r.Name,
			CourseID: r.CourseID,
			Course: &generated.Course{
				ID:        r.Course.ID,
				Name:      r.Course.Name,
				UpdatedAt: r.Course.UpdatedAt,
				CreatedAt: r.Course.CreatedAt,
			},
			UpdatedAt: r.UpdatedAt,
			CreatedAt: r.CreatedAt,
		}}, nil
	}

	skip, take := getPaginationBounds(pagination)

	classesList, err := r.server.educationClient.GetClasses(ctx, skip, take)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var classes []*generated.Class
	for _, class := range classesList {
		classes = append(classes, &generated.Class{
			ID:       class.ID,
			Name:     class.Name,
			CourseID: class.CourseID,
			Course: &generated.Course{
				ID:        class.Course.ID,
				Name:      class.Course.Name,
				UpdatedAt: class.Course.UpdatedAt,
				CreatedAt: class.Course.CreatedAt,
			},
			UpdatedAt: class.UpdatedAt,
			CreatedAt: class.CreatedAt,
		})
	}
	return classes, nil
}

func getPaginationBounds(p *generated.PaginationInput) (uint64, uint64) {
	skipValue := uint64(0)
	takeValue := uint64(100)
	if p != nil {
		if p.Skip != nil {
			skipValue = uint64(*p.Skip)
		}
		if p.Take != nil {
			takeValue = uint64(*p.Take)
		}
	}
	return skipValue, takeValue
}
