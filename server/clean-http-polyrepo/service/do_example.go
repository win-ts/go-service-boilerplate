package service

import (
	"context"
)

// DoExample returns example string
func (s *service) DoExample(ctx context.Context) (string, error) {
	res, err := s.exampleRepository.DoExample(ctx)
	if err != nil {
		return "", err
	}

	return res, nil
}
