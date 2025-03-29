package course

import (
	"context"
	"errors"
	"fmt"
	"github.com/og11423074s/gocourse_domain/domain"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

type (
	Repository interface {
		Create(ctx context.Context, course *domain.Course) error
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error)
		Get(ctx context.Context, id string) (*domain.Course, error)
		DeleteById(ctx context.Context, id string) error
		Update(ctx context.Context, id string, name *string, startDate, endDate *time.Time) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	repo struct {
		db  *gorm.DB
		log *log.Logger
	}
)

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		db:  db,
		log: log,
	}
}

func (r *repo) Create(ctx context.Context, course *domain.Course) error {

	if err := r.db.WithContext(ctx).Create(course); err.Error != nil {
		r.log.Println(err.Error)
		return err.Error
	}
	r.log.Println("Course created with id: ", course.ID)
	return nil
}

func (r *repo) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error) {
	var courses []domain.Course

	tx := r.db.WithContext(ctx).Model(&courses)
	tx = applyFilters(tx, filters)
	tx = tx.Offset(offset).Limit(limit)

	result := tx.Order("created_at desc").Find(&courses)

	if result.Error != nil {
		r.log.Println(result.Error)
		return nil, result.Error
	}

	return courses, nil
}

func (r *repo) Get(ctx context.Context, id string) (*domain.Course, error) {
	course := domain.Course{ID: id}

	err := r.db.WithContext(ctx).First(&course).Error
	if err != nil {
		r.log.Println(err)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorNotFound{id}
		}

		return nil, err
	}

	return &course, nil
}

func (r *repo) DeleteById(ctx context.Context, id string) error {
	user := domain.Course{ID: id}

	result := r.db.WithContext(ctx).Delete(&user)

	if result.Error != nil {
		r.log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.log.Printf("course %s doesn't exists", id)
		return ErrorNotFound{id}
	}

	return nil
}

func (r *repo) Update(ctx context.Context, id string, name *string, startDate, endDate *time.Time) error {

	values := make(map[string]interface{})

	if name != nil {
		values["name"] = *name
	}

	if !startDate.IsZero() {
		values["start_date"] = startDate

	}

	if !endDate.IsZero() {
		values["end_date"] = endDate

	}

	result := r.db.WithContext(ctx).Model(&domain.Course{}).Where("id = ?", id).UpdateColumns(values)

	if result.Error != nil {
		r.log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.log.Printf("course %s doesn't exists", id)
		return ErrorNotFound{id}
	}

	return nil

}

func (r *repo) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	tx := r.db.WithContext(ctx).Model(&domain.Course{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		r.log.Println(err)
		return 0, err
	}

	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.name != "" {
		filters.name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.name))
		tx = tx.Where("lower(name) like ?", filters.name)
	}

	return tx
}
