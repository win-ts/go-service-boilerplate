package repository

import (
	"context"
)

// List returns a list of example entities from database
func (r *exampleRepository) DoExample(ctx context.Context) (string, error) {
	_ = ctx
	return "example", nil
}
