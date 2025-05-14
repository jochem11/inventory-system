package education

import (
	"context"
	"github.com/segmentio/ksuid"
	"time"
)

type Service interface {
	PostCourse(ctx context.Context, name string) (*Course, error)
	GetCourse(ctx context.Context, id string) (*Course, error)
	GetCourses(ctx context.Context, skip *uint64, take *uint64) ([]*Course, error)
	DeleteCourseByID(ctx context.Context, id string) error
	UpdateCourse(ctx context.Context, id string, name *string) (*Course, error)

	PostClass(ctx context.Context, name, courseID string) (*Class, error)
	GetClass(ctx context.Context, id string) (*Class, error)
	GetClasses(ctx context.Context, skip *uint64, take *uint64) ([]*Class, error)
	DeleteClassByID(ctx context.Context, id string) error
	UpdateClass(ctx context.Context, id string, name *string, courseID *string) (*Class, error)
}

func NewEducationService(r Repository) Service {
	return &educationService{r}
}

type Course struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Class struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CourseID  string    `json:"course_id"`
	Course    *Course   `json:"course,omitempty"` // optional: populated when joined
}

type educationService struct {
	repository Repository
}

func (s *educationService) defaultSkipTake(skip *uint64, take *uint64) (*uint64, *uint64) {
	const defaultSkip uint64 = 0
	const defaultTake uint64 = 50

	if skip == nil {
		skip = new(uint64)
		*skip = defaultSkip
	}
	if take == nil {
		take = new(uint64)
		*take = defaultTake
	}

	return skip, take
}

func (s *educationService) PostCourse(ctx context.Context, name string) (*Course, error) {
	c := &Course{
		ID:        ksuid.New().String(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repository.PutCourse(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *educationService) GetCourse(ctx context.Context, id string) (*Course, error) {
	return s.repository.GetCourseByID(ctx, id)
}

func (s *educationService) GetCourses(ctx context.Context, skip *uint64, take *uint64) ([]*Course, error) {
	skip, take = s.defaultSkipTake(skip, take)
	return s.repository.ListCourses(ctx, *skip, *take)
}

func (s *educationService) DeleteCourseByID(ctx context.Context, id string) error {
	return s.repository.DeleteCourseByID(ctx, id)
}

func (s *educationService) UpdateCourse(ctx context.Context, id string, name *string) (*Course, error) {
	existing, err := s.repository.GetCourseByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != nil {
		existing.Name = *name
	}

	existing.UpdatedAt = time.Now()
	if err := s.repository.PutCourse(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *educationService) PostClass(ctx context.Context, name, courseID string) (*Class, error) {
	c := &Class{
		ID:        ksuid.New().String(),
		Name:      name,
		CourseID:  courseID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repository.PutClass(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *educationService) GetClass(ctx context.Context, id string) (*Class, error) {
	return s.GetClass(ctx, id)
}

func (s *educationService) GetClasses(ctx context.Context, skip *uint64, take *uint64) ([]*Class, error) {
	skip, take = s.defaultSkipTake(skip, take)
	return s.repository.ListClasses(ctx, *skip, *take)
}

func (s *educationService) DeleteClassByID(ctx context.Context, id string) error {
	return s.repository.DeleteClassByID(ctx, id)
}

func (s *educationService) UpdateClass(ctx context.Context, id string, name *string, courseID *string) (*Class, error) {
	existing, err := s.repository.GetClassByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != nil {
		existing.Name = *name
	}
	if courseID != nil {
		existing.CourseID = *courseID
	}

	existing.UpdatedAt = time.Now()

	if err := s.repository.PutClass(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}
