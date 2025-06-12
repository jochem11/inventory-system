package education

import (
	"context"
	"github.com/jochem11/inventory-system-back/education/pb"
	"google.golang.org/grpc"
	"io"
	"log"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.EducationServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := pb.NewEducationServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostCourse(ctx context.Context, name string) (*Course, error) {
	r, err := c.service.PostCourse(ctx, &pb.PostCourseRequest{
		Name: name,
	})

	if err != nil {
		return nil, err
	}

	updatedAt := r.Course.UpdatedAt.AsTime()
	createdAt := r.Course.CreatedAt.AsTime()

	return &Course{
		ID:        r.Course.Id,
		Name:      r.Course.Name,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
	}, nil
}

func (c *Client) GetCourse(ctx context.Context, id string) (*Course, error) {
	r, err := c.service.GetCourse(ctx, &pb.GetCourseRequest{Id: id})
	if err != nil {
		return nil, err
	}

	updatedAt := r.Course.UpdatedAt.AsTime()
	createdAt := r.Course.CreatedAt.AsTime()

	return &Course{
		ID:        r.Course.Id,
		Name:      r.Course.Name,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
	}, nil
}

func (c *Client) GetCourses(ctx context.Context, skip, take uint64) ([]*Course, error) {
	r, err := c.service.GetCourses(ctx, &pb.GetCoursesRequest{Skip: skip, Take: take})
	if err != nil {
		return nil, err
	}

	courses := []*Course{}
	for _, course := range r.Courses {
		updatedAt := course.UpdatedAt.AsTime()
		createdAt := course.CreatedAt.AsTime()
		courses = append(courses, &Course{
			ID:        course.Id,
			Name:      course.Name,
			UpdatedAt: updatedAt,
			CreatedAt: createdAt,
		})
	}

	return courses, nil
}

func (c *Client) UpdateCourse(ctx context.Context, id string, name *string) (*Course, error) {
	r, err := c.service.UpdateCourse(ctx, &pb.UpdateCourseRequest{Id: id, Name: name})
	if err != nil {
		return nil, err
	}
	updatedAt := r.Course.UpdatedAt.AsTime()
	createdAt := r.Course.CreatedAt.AsTime()

	return &Course{
		ID:        r.Course.Id,
		Name:      r.Course.Name,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
	}, nil
}

func (c *Client) DeleteCourse(ctx context.Context, id string) error {
	_, err := c.service.DeleteCourse(ctx, &pb.DeleteCourseRequest{Id: id})
	return err
}

func (c *Client) LiveCourses(ctx context.Context, skip, take uint64) (<-chan []*Course, error) {
	stream, err := c.service.LiveCourses(ctx, &pb.GetCoursesRequest{
		Skip: skip,
		Take: take,
	})
	if err != nil {
		return nil, err
	}

	ch := make(chan []*Course)
	go func() {
		defer close(ch)
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				// Log the error but don't send it through the channel
				log.Printf("Error receiving stream: %v", err)
				return
			}

			courses := make([]*Course, 0, len(resp.Courses))
			for _, c := range resp.Courses {
				courses = append(courses, &Course{
					ID:        c.Id,
					Name:      c.Name,
					UpdatedAt: c.UpdatedAt.AsTime(),
					CreatedAt: c.CreatedAt.AsTime(),
				})
			}

			select {
			case ch <- courses:
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch, nil
}

func (c *Client) PostClass(ctx context.Context, courseId, name string) (*Class, error) {
	r, err := c.service.PostClass(ctx, &pb.PostClassRequest{CourseId: courseId, Name: name})
	if err != nil {
		return nil, err
	}

	updatedAt := r.Class.UpdatedAt.AsTime()
	createdAt := r.Class.CreatedAt.AsTime()

	courseUpdatedAt := r.Class.Course.UpdatedAt.AsTime()
	courseCreatedAt := r.Class.Course.CreatedAt.AsTime()

	return &Class{
		ID:       r.Class.Id,
		Name:     r.Class.Name,
		CourseID: r.Class.CourseId,
		Course: &Course{
			ID:        r.Class.Course.Id,
			Name:      r.Class.Course.Name,
			UpdatedAt: courseUpdatedAt,
			CreatedAt: courseCreatedAt,
		},
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
	}, nil
}

func (c *Client) GetClass(ctx context.Context, id string) (*Class, error) {
	r, err := c.service.GetClass(ctx, &pb.GetClassRequest{Id: id})
	if err != nil {
		return nil, err
	}

	updatedAt := r.Class.UpdatedAt.AsTime()
	createdAt := r.Class.CreatedAt.AsTime()

	courseUpdatedAt := r.Class.Course.UpdatedAt.AsTime()
	courseCreatedAt := r.Class.Course.CreatedAt.AsTime()

	return &Class{
		ID:       r.Class.Id,
		Name:     r.Class.Name,
		CourseID: r.Class.CourseId,
		Course: &Course{
			ID:        r.Class.Course.Id,
			Name:      r.Class.Course.Name,
			UpdatedAt: courseUpdatedAt,
			CreatedAt: courseCreatedAt,
		},
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
	}, nil
}

func (c *Client) GetClasses(ctx context.Context, skip, take uint64) ([]*Class, error) {
	r, err := c.service.GetClasses(ctx, &pb.GetClassesRequest{Skip: skip, Take: take})
	if err != nil {
		return nil, err
	}

	classes := []*Class{}
	for _, class := range r.Classes {
		updatedAt := class.UpdatedAt.AsTime()
		createdAt := class.CreatedAt.AsTime()

		courseUpdatedAt := class.Course.UpdatedAt.AsTime()
		courseCreatedAt := class.Course.CreatedAt.AsTime()
		classes = append(classes, &Class{
			ID:       class.Id,
			Name:     class.Name,
			CourseID: class.CourseId,
			Course: &Course{
				ID:        class.Course.Id,
				Name:      class.Course.Name,
				UpdatedAt: courseUpdatedAt,
				CreatedAt: courseCreatedAt,
			},
			UpdatedAt: updatedAt,
			CreatedAt: createdAt,
		})
	}
	return classes, nil
}

func (c *Client) UpdateClass(ctx context.Context, id string, name, courseId *string) (*Class, error) {
	r, err := c.service.UpdateClass(ctx, &pb.UpdateClassRequest{Id: id, Name: name, CourseId: courseId})
	if err != nil {
		return nil, err
	}

	updatedAt := r.Class.UpdatedAt.AsTime()
	createdAt := r.Class.CreatedAt.AsTime()

	courseUpdatedAt := r.Class.Course.UpdatedAt.AsTime()
	courseCreatedAt := r.Class.Course.CreatedAt.AsTime()

	return &Class{
		ID:       r.Class.Id,
		Name:     r.Class.Name,
		CourseID: r.Class.CourseId,
		Course: &Course{
			ID:        r.Class.Course.Id,
			Name:      r.Class.Course.Name,
			UpdatedAt: courseUpdatedAt,
			CreatedAt: courseCreatedAt,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (c *Client) DeleteClass(ctx context.Context, id string) error {
	_, err := c.service.DeleteClass(ctx, &pb.DeleteClassRequest{Id: id})
	return err
}
