package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Informasjonsforvaltning/catalog-history-service/config/postgresql"
	"github.com/Informasjonsforvaltning/catalog-history-service/logging"
	"github.com/Informasjonsforvaltning/catalog-history-service/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type UpdateRepository interface {
	StoreUpdate(ctx context.Context, update model.Update) error
	GetUpdates(ctx context.Context, catalogId string, resourceId *string, page int, size int, sortBy string, sortOrder string) ([]model.Update, int64, error)
	GetUpdate(ctx context.Context, catalogId string, resourceId string, updateId string) (*model.Update, error)
}

type UpdateRepositoryImpl struct {
	pool *pgxpool.Pool
}

var updateRepository *UpdateRepositoryImpl

func InitRepository() *UpdateRepositoryImpl {
	if updateRepository == nil {
		updateRepository = &UpdateRepositoryImpl{pool: postgresql.Pool()}
	}
	return updateRepository
}

func (r UpdateRepositoryImpl) StoreUpdate(ctx context.Context, update model.Update) error {
	opsJSON, err := json.Marshal(update.Operations)
	if err != nil {
		logrus.Errorf("failed to marshal operations: %v", err)
		logging.LogAndPrintError(err)
		return err
	}

	_, err = r.pool.Exec(ctx,
		`INSERT INTO updates (id, catalog_id, resource_id, person_id, person_email, person_name, datetime, operations)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		update.ID,
		update.CatalogId,
		update.ResourceId,
		update.Person.ID,
		update.Person.Email,
		update.Person.Name,
		update.DateTime,
		opsJSON,
	)
	if err != nil {
		logging.LogAndPrintError(err)
		return err
	}
	return nil
}

func (r UpdateRepositoryImpl) GetUpdates(ctx context.Context, catalogId string, resourceId *string, page int, size int, sortBy string, sortOrder string) ([]model.Update, int64, error) {
	validatedPage, validatedSize, err := ValidatePagination(page, size)
	if err != nil {
		return nil, 0, err
	}
	page = validatedPage
	size = validatedSize

	validatedSortBy := ValidateSortField(sortBy)

	offset := page * size

	var whereClause string
	var args []any
	if resourceId != nil {
		whereClause = "WHERE catalog_id = $1 AND resource_id = $2"
		args = []any{catalogId, *resourceId}
	} else {
		whereClause = "WHERE catalog_id = $1"
		args = []any{catalogId}
	}

	query := fmt.Sprintf(
		`SELECT id, catalog_id, resource_id, person_id, person_email, person_name, datetime, operations
		 FROM updates %s ORDER BY %s %s LIMIT %d OFFSET %d`,
		whereClause, validatedSortBy, sortOrder, size, offset,
	)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		logrus.Errorf("query failed: %v", err)
		logging.LogAndPrintError(err)
		return nil, 0, err
	}
	defer rows.Close()

	var updates []model.Update
	for rows.Next() {
		var u model.Update
		var opsJSON []byte
		var dt *time.Time
		err := rows.Scan(
			&u.ID, &u.CatalogId, &u.ResourceId,
			&u.Person.ID, &u.Person.Email, &u.Person.Name,
			&dt, &opsJSON,
		)
		if dt != nil {
			utc := dt.UTC()
			u.DateTime = &utc
		}
		if err != nil {
			logrus.Errorf("scan failed: %v", err)
			logging.LogAndPrintError(err)
			return nil, 0, err
		}
		if err := json.Unmarshal(opsJSON, &u.Operations); err != nil {
			logrus.Errorf("unmarshal operations failed: %v", err)
			logging.LogAndPrintError(err)
			return nil, 0, err
		}
		updates = append(updates, u)
	}
	if err := rows.Err(); err != nil {
		logging.LogAndPrintError(err)
		return nil, 0, err
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM updates %s", whereClause)
	var count int64
	err = r.pool.QueryRow(ctx, countQuery, args...).Scan(&count)
	if err != nil {
		logrus.Errorf("count query failed: %v", err)
		logging.LogAndPrintError(err)
		return nil, 0, err
	}

	return updates, count, nil
}

func (r UpdateRepositoryImpl) GetUpdate(ctx context.Context, catalogId string, resourceId string, updateId string) (*model.Update, error) {
	logrus.Info("Starting to get update from database")

	if err := ValidateID(catalogId, "catalogId"); err != nil {
		logrus.Errorf("Invalid catalogId: %v", err)
		return nil, err
	}
	if err := ValidateID(resourceId, "resourceId"); err != nil {
		logrus.Errorf("Invalid resourceId: %v", err)
		return nil, err
	}
	if err := ValidateID(updateId, "updateId"); err != nil {
		logrus.Errorf("Invalid updateId: %v", err)
		return nil, err
	}

	var u model.Update
	var opsJSON []byte
	var dt *time.Time
	err := r.pool.QueryRow(ctx,
		`SELECT id, catalog_id, resource_id, person_id, person_email, person_name, datetime, operations
		 FROM updates WHERE id = $1 AND catalog_id = $2 AND resource_id = $3`,
		updateId, catalogId, resourceId,
	).Scan(
		&u.ID, &u.CatalogId, &u.ResourceId,
		&u.Person.ID, &u.Person.Email, &u.Person.Name,
		&dt, &opsJSON,
	)
	if dt != nil {
		utc := dt.UTC()
		u.DateTime = &utc
	}
	if err == pgx.ErrNoRows {
		logrus.Error("update not found in db")
		return nil, nil
	}
	if err != nil {
		logrus.Errorf("error when getting update from db: %s", err)
		logging.LogAndPrintError(err)
		return nil, err
	}

	if err := json.Unmarshal(opsJSON, &u.Operations); err != nil {
		logrus.Errorf("error when unmarshalling operations: %s", err)
		logging.LogAndPrintError(err)
		return nil, err
	}

	return &u, nil
}
