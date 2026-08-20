package auth

import (
	"context"
	"strings"
	"time"

	"omniforge-api/internal/cache"
	"omniforge-api/internal/role"
	"omniforge-api/internal/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	userRepository *user.Repository
	roleRepository *role.Repository
	redisClient     *redis.Client
	jwtSecret       string
}

func NewService(
	userRepository *user.Repository,
	roleRepository *role.Repository,
	redisClient *redis.Client,
	jwtSecret string,
) *Service {
	return &Service{
		userRepository: userRepository,
		roleRepository: roleRepository,
		redisClient:     redisClient,
		jwtSecret:       jwtSecret,
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

func (s *Service) Login(
	ctx context.Context,
	input LoginRequest,
) (*LoginResult, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	existingUser, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, ErrInvalidCredentials
	}

	if err := ComparePassword(
		existingUser.PasswordHash,
		input.Password,
	); err != nil {
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

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	err = cache.StoreRefreshToken(
		ctx,
		s.redisClient,
		refreshToken,
		cache.RefreshSession{
			UserID: existingUser.ID,
			Role:   existingUser.Role.Name,
		},
		RefreshTokenTTL,
	)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		ID:           existingUser.ID,
		Name:         existingUser.Name,
		Email:        existingUser.Email,
		Role:         existingUser.Role.Name,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Logout(
	ctx context.Context,
	accessToken string,
	refreshToken string,
) error {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		accessToken,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(s.jwtSecret), nil
		},
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Alg(),
		}),
		jwt.WithExpirationRequired(),
	)

	if err != nil || !token.Valid {
		return ErrInvalidToken
	}

	if claims.ExpiresAt == nil {
		return ErrInvalidToken
	}

	ttl := time.Until(claims.ExpiresAt.Time)

	if ttl <= 0 {
		return ErrInvalidToken
	}

	if err := cache.DeleteRefreshToken(
		ctx,
		s.redisClient,
		refreshToken,
	); err != nil {
		return err
	}

	if err := cache.BlacklistToken(
		ctx,
		s.redisClient,
		accessToken,
		ttl,
	); err != nil {
		return err
	}

	return nil
}
func (s *Service) Refresh(
	ctx context.Context,
	refreshToken string,
) (*RefreshResult, error) {
	session, found, err := cache.GetRefreshSession(
		ctx,
		s.redisClient,
		refreshToken,
	)

	if err != nil {
		return nil, err
	}

	if !found {
		return nil, ErrInvalidRefreshToken
	}

	accessToken, err := GenerateAccessToken(
		session.UserID,
		session.Role,
		s.jwtSecret,
	)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	err = cache.StoreRefreshToken(
		ctx,
		s.redisClient,
		newRefreshToken,
		*session,
		RefreshTokenTTL,
	)
	if err != nil {
		return nil, err
	}

	if err := cache.DeleteRefreshToken(
		ctx,
		s.redisClient,
		refreshToken,
	); err != nil {
		_ = cache.DeleteRefreshToken(
			ctx,
			s.redisClient,
			newRefreshToken,
		)

		return nil, err
	}

	return &RefreshResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}