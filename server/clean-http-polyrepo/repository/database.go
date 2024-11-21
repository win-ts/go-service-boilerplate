package repository

import (
	"database/sql"

	"github.com/win-ts/go-service-boilerplate/server/clean-http-polyrepo/dto"
)

type databaseRepository struct {
	database string
	client   *sql.DB
}

// DatabaseRepositoryConfig represents the configuration for wiremock API repository
type DatabaseRepositoryConfig struct {
	Database string
}

// DatabaseRepositoryDependencies represents the dependencies for wiremock API repository
type DatabaseRepositoryDependencies struct {
	Client *sql.DB
}

// NewDatabaseRepository creates a new wiremock API repository
func NewDatabaseRepository(c DatabaseRepositoryConfig, d DatabaseRepositoryDependencies) DatabaseRepository {
	return &databaseRepository{
		database: c.Database,
		client:   d.Client,
	}
}

// QueryTest returns the test rows queried from test table
func (r *databaseRepository) QueryTest() (*[]dto.TestEntity, error) {
	tests := []dto.TestEntity{}

	query := `
		SELECT * FROM tbl_test;
	`

	rows, err := r.client.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			return
		}
	}()

	for rows.Next() {
		test := dto.TestEntity{}
		if err := rows.Scan(&test.ID, &test.Message); err != nil {
			return nil, err
		}
		tests = append(tests, test)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &tests, nil
}
