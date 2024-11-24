package service

import "github.com/win-ts/go-service-boilerplate/worker/clean-worker/dto"

func (s *service) DoDBTest() (*[]dto.TestModel, error) {
	entity, err := s.databaseRepository.QueryTest()
	if err != nil {
		return nil, err
	}

	res := []dto.TestModel{}
	for _, e := range *entity {
		res = append(res, dto.TestModel(e))
	}

	return &res, nil
}
