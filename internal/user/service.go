package user

type Service struct {
	userRepository *Repository
}

func NewService(userRepository *Repository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

func (s *Service) GetProfile(userID uint) (*Response, error) {
	existingUser, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, ErrUserNotFound
	}

	result := ToResponse(existingUser)

	return &result, nil
}