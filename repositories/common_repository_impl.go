package repositories

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/satryarangga/amartha-loan-engine/models"
	"gorm.io/gorm"
)

type CommonRepositoryImpl[T any] struct {
	db *gorm.DB
}

func NewCommonRepository[T any](db *gorm.DB) *CommonRepositoryImpl[T] {
	return &CommonRepositoryImpl[T]{db: db}
}

func (r *CommonRepositoryImpl[T]) buildQueryFindAll(param models.FindAllParam, query *gorm.DB) *gorm.DB {
	if param.SortBy.FieldName != "" {
		query = query.Order(fmt.Sprintf("%s %s", param.SortBy.FieldName, param.SortBy.Direction))
	}

	if param.Limit > 0 {
		query = query.
			Limit(int(param.Limit)).
			Offset((param.Offset - 1) * param.Limit)
	}

	if param.SearchKeyword != "" && len(param.FieldsToSearch) > 0 {
		conditions := make([]string, 0, len(param.FieldsToSearch))
		args := make([]interface{}, 0, len(param.FieldsToSearch))

		for _, field := range param.FieldsToSearch {
			conditions = append(conditions, field+" ILIKE ?")
			args = append(args, "%"+param.SearchKeyword+"%")
		}

		whereClause := "(" + strings.Join(conditions, " OR ") + ")"
		query = query.Where(whereClause, args...)
	}

	for _, relation := range param.PreloadTables {
		query = query.Preload(relation)
	}

	for _, relation := range param.JoinTables {
		query = query.Joins(relation)
	}

	return query
}

func (r *CommonRepositoryImpl[T]) FindAll(ctx context.Context, param models.FindAllParam) ([]T, error) {
	var models []T
	query := r.db.WithContext(ctx)
	query = r.buildQueryFindAll(param, query)

	// Execute the query
	result := query.Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}
	return models, nil
}

func (r *CommonRepositoryImpl[T]) Insert(ctx context.Context, tx *gorm.DB, model *T) (string, error) {
	db := r.db
	if tx != nil {
		db = tx
	}

	result := db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return "", result.Error
	}

	idField := reflect.ValueOf(model).Elem().FieldByName("ID")
	if !idField.IsValid() {
		return "", nil
	}
	id := idField.String()
	return id, nil
}

func (r *CommonRepositoryImpl[T]) Update(ctx context.Context, tx *gorm.DB, model *T) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	result := db.WithContext(ctx).Save(model)
	return result.Error
}

func (r *CommonRepositoryImpl[T]) FindByID(ctx context.Context, id string) (*T, error) {
	var model T
	query := r.db.WithContext(ctx).Where("id = ?", id)

	result := query.First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return &model, nil
}
