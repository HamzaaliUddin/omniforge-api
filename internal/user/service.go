package user

import (
	"context"
	"strings"

	"omniforge-api/internal/cache"
	"omniforge-api/internal/role"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	repository  *Repository
	redisClient *redis.Client
}

func NewService(
	repository *Repository,
	redisClient *redis.Client,
) *Service {
	return &Service{
		repository:  repository,
		redisClient: redisClient,
	}
}

func (s *Service) GetMe(
	ctx context.Context,
	userID uint,
) (*cache.UserCache, error) {
	cachedUser, found, err := cache.GetUser(
		ctx,
		s.redisClient,
		userID,
	)

	if err == nil && found {
		return cachedUser, nil
	}

	user, err := s.repository.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	result := &cache.UserCache{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role.Name,
	}

	_ = cache.SetUser(
		ctx,
		s.redisClient,
		*result,
	)

	return result, nil
}

func (s *Service) UpdateMe(
	ctx context.Context,
	userID uint,
	input UpdateMeRequest,
) (*UserResponse, error) {
	existingUser, err := s.repository.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, ErrUserNotFound
	}

	if input.Name != nil {
		existingUser.Name = strings.TrimSpace(
			*input.Name,
		)
	}

	if input.Email != nil {
		email := strings.ToLower(
			strings.TrimSpace(*input.Email),
		)

		if email != existingUser.Email {
			userWithEmail, err := s.repository.FindByEmail(
				email,
			)
			if err != nil {
				return nil, err
			}

			if userWithEmail != nil &&
				userWithEmail.ID != userID {
				return nil, ErrEmailAlreadyExists
			}

			existingUser.Email = email
		}
	}

	if err := s.repository.Update(existingUser); err != nil {
		return nil, err
	}

	_ = cache.DeleteUser(
		ctx,
		s.redisClient,
		userID,
	)

	result := ToResponse(existingUser)

	return &result, nil
}

func (s *Service) DeleteMe(
	ctx context.Context,
	userID uint,
) error {
	existingUser, err := s.repository.FindByID(userID)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return ErrUserNotFound
	}

	if existingUser.Role.Name == role.NameAdmin {
	return ErrAdminCannotDeleteSelf
	}

	if err := s.repository.Delete(userID); err != nil {
		return err
	}

	_ = cache.DeleteUser(
		ctx,
		s.redisClient,
		userID,
	)

	return nil
}