package education

import (
	"context"
	"database/sql"
)

type Repository interface {
	Close()

	PutCourse(ctx context.Context, c *Course) error
	GetCourseByID(ctx context.Context, id string) (*Course, error)
	ListCourses(ctx context.Context, skip uint64, take uint64) ([]*Course, error)
	UpdateCourse(ctx context.Context, c *Course) (*Course, error)
	DeleteCourseByID(ctx context.Context, id string) error

	PutClass(ctx context.Context, c *Class) error
	GetClassByID(ctx context.Context, id string) (*Class, error)
	ListClasses(ctx context.Context, skip uint64, take uint64) ([]*Class, error)
	UpdateClass(ctx context.Context, c *Class) (*Class, error)
	DeleteClassByID(ctx context.Context, id string) error
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db}, nil
}

func (r postgresRepository) Close() {
	r.db.Close()

}

func (r *postgresRepository) PutCourse(ctx context.Context, c *Course) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO courses(id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)", c.ID, c.Name, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *postgresRepository) GetCourseByID(ctx context.Context, id string) (*Course, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM courses WHERE id = $1", id)
	c := &Course{}
	if err := row.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	return c, nil
}

func (r *postgresRepository) ListCourses(ctx context.Context, skip uint64, take uint64) ([]*Course, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, created_at, updated_at FROM courses ORDER BY id DESC OFFSET $1 LIMIT $2",
		skip,
		take,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	courses := []*Course{}
	for rows.Next() {
		c := &Course{}
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *postgresRepository) UpdateCourse(ctx context.Context, c *Course) (*Course, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE courses SET name = $1, updated_at = $2 WHERE id = $3", c.Name, c.UpdatedAt, c.ID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return r.GetCourseByID(ctx, c.ID)
}

func (r *postgresRepository) DeleteCourseByID(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM courses WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *postgresRepository) PutClass(ctx context.Context, c *Class) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO classes(id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)", c.ID, c.Name, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *postgresRepository) GetClassByID(ctx context.Context, id string) (*Class, error) {
	row := r.db.QueryRowContext(ctx, `
        SELECT c.id, c.name, c.created_at, c.updated_at, 
               cl.id, cl.name, cl.created_at, cl.updated_at
        FROM classes cl
        JOIN courses c ON cl.course_id = c.id
        WHERE cl.id = $1`, id)

	class := &Class{}
	course := &Course{}
	err := row.Scan(&course.ID, &course.Name, &course.CreatedAt, &course.UpdatedAt,
		&class.ID, &class.Name, &class.CreatedAt, &class.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Associate the course with the class
	class.Course = course

	return class, nil
}

func (r *postgresRepository) ListClasses(ctx context.Context, skip uint64, take uint64) ([]*Class, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT c.id, c.name, c.created_at, c.updated_at, cl.id, cl.name, cl.created_at, cl.updated_at
        FROM classes cl
        JOIN courses c ON cl.course_id = c.id
        ORDER BY cl.id DESC
        OFFSET $1 LIMIT $2`, skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	classes := []*Class{}
	for rows.Next() {
		course := &Course{}
		class := &Class{}
		if err := rows.Scan(&course.ID, &course.Name, &course.CreatedAt, &course.UpdatedAt,
			&class.ID, &class.Name, &class.CreatedAt, &class.UpdatedAt); err != nil {
			return nil, err
		}
		// Associate the course with the class
		class.Course = course
		classes = append(classes, class)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return classes, nil
}

func (r *postgresRepository) UpdateClass(ctx context.Context, c *Class) (*Class, error) {
	_, err := r.db.ExecContext(ctx, `
        UPDATE classes 
        SET name = $1, updated_at = $2 
        WHERE id = $3`, c.Name, c.UpdatedAt, c.ID)
	if err != nil {
		return nil, err
	}

	return r.GetClassByID(ctx, c.ID)
}

func (r *postgresRepository) DeleteClassByID(ctx context.Context, id string) error {
	// Delete the class (the associated course will be automatically deleted due to ON DELETE CASCADE)
	res, err := r.db.ExecContext(ctx, "DELETE FROM classes WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Class not found
	}

	return nil
}
