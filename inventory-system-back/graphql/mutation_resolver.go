package main

import (
	"context"
	"github.com/jochem11/inventory-system-back/graphql/generated"
	"log"
	"time"
)

type mutationResolver struct {
	server *Server
}

// Courses
func (r mutationResolver) CreateCourse(ctx context.Context, course generated.CreateCourseInput) (*generated.Course, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	c, err := r.server.educationClient.PostCourse(ctx, course.Name)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &generated.Course{
		ID:        c.ID,
		Name:      c.Name,
		UpdatedAt: c.UpdatedAt,
		CreatedAt: c.CreatedAt,
	}, nil
}

func (r mutationResolver) UpdateCourse(ctx context.Context, course generated.UpdateCourseInput) (*generated.Course, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	c, err := r.server.educationClient.UpdateCourse(ctx, course.ID, course.Name)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &generated.Course{
		ID:        c.ID,
		Name:      c.Name,
		UpdatedAt: c.UpdatedAt,
		CreatedAt: c.CreatedAt,
	}, nil
}

func (r mutationResolver) DeleteCourse(ctx context.Context, course generated.DeleteByIDCourseInput) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.server.educationClient.DeleteCourse(ctx, course.ID)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}

// Classes
func (r mutationResolver) CreateClass(ctx context.Context, class generated.CreateClassInput) (*generated.Class, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	c, err := r.server.educationClient.PostClass(ctx, class.Name, class.CourseID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &generated.Class{
		ID:       c.ID,
		Name:     c.Name,
		CourseID: c.CourseID,
		Course: &generated.Course{
			ID:        c.Course.ID,
			Name:      c.Course.Name,
			UpdatedAt: c.Course.UpdatedAt,
			CreatedAt: c.Course.CreatedAt,
		},
		UpdatedAt: c.UpdatedAt,
		CreatedAt: c.CreatedAt,
	}, nil
}

func (r mutationResolver) UpdateClass(ctx context.Context, class generated.UpdateClassInput) (*generated.Class, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	c, err := r.server.educationClient.UpdateClass(ctx, class.ID, class.Name, class.CourseID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &generated.Class{
		ID:       c.ID,
		Name:     c.Name,
		CourseID: c.CourseID,
		Course: &generated.Course{
			ID:        c.Course.ID,
			Name:      c.Course.Name,
			UpdatedAt: c.Course.UpdatedAt,
			CreatedAt: c.Course.CreatedAt,
		},
		UpdatedAt: c.UpdatedAt,
		CreatedAt: c.CreatedAt,
	}, nil
}

func (r mutationResolver) DeleteClass(ctx context.Context, class generated.DeleteByIDClassInput) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.server.educationClient.DeleteClass(ctx, class.ID)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}
