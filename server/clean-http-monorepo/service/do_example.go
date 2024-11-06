package service

import (
	"context"
)

// DoExample returns example string
func (s *service) DoExample(pctx context.Context) (string, error) {
	res, err := s.exampleRepository.DoExample(pctx)
	if err != nil {
		return "", err
	}

	return res, nil
}
