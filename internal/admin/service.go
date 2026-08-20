package admin

import "omniforge-api/internal/user"

type Service struct {
	userRepository *user.Repository
}

func NewService(userRepository *user.Repository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

func (s *Service) GetUsers() ([]user.UserResponse, error) {
	users, err := s.userRepository.FindAll()
	if err != nil {
		return nil, err
	}

	result := make([]user.UserResponse, 0, len(users))

	for i := range users {
		result = append(
			result,
			user.ToResponse(&users[i]),
		)
	}

	return result, nil
}