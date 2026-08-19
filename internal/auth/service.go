package auth

import (
	"strings"

	"omniforge-api/internal/role"
	"omniforge-api/internal/user"
)

type Service struct {
	userRepository *user.Repository
	roleRepository *role.Repository
	jwtSecret      string
}

func NewService(
	userRepository *user.Repository,
	roleRepository *role.Repository,
	jwtSecret string,
) *Service {
	return &Service{
		userRepository: userRepository,
		roleRepository: roleRepository,
		jwtSecret:      jwtSecret,
	}
}

func (s *Service) Register(input RegisterRequest) (*user.User, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	existingUser, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	defaultRole, err := s.roleRepository.FindByName("user")
	if err != nil {
		return nil, err
	}

	if defaultRole == nil {
		return nil, ErrDefaultRoleNotFound
	}

	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	newUser := &user.User{
		Name:         strings.TrimSpace(input.Name),
		Email:        email,
		PasswordHash: passwordHash,
		RoleID:       defaultRole.ID,
		Role:         *defaultRole,
	}

	if err := s.userRepository.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *Service) Login(input LoginRequest) (*LoginResult, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	existingUser, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, ErrInvalidCredentials
	}

	if err := ComparePassword(existingUser.PasswordHash, input.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := GenerateAccessToken(
		existingUser.ID,
		existingUser.Role.Name,
		s.jwtSecret,
	)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		ID:          existingUser.ID,
		Name:        existingUser.Name,
		Email:       existingUser.Email,
		Role:        existingUser.Role.Name,
		AccessToken: accessToken,
	}, nil
}