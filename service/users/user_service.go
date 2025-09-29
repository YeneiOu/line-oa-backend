package users

import (
	"context"
	"errors"
	"line-oa-backend/entities"
	"line-oa-backend/repositories/users"
)

type UserService struct {
	userRepo *users.UserRepository
}

func NewUserService(userRepo *users.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, lineUserID, name, email, pictureURL string) (*entities.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByLineUserID(ctx, lineUserID)
	if err == nil && existingUser != nil {
		return existingUser, nil
	}

	// Create new user
	user := entities.NewUser(lineUserID, name, email, pictureURL)
	if !user.IsValid() {
		return nil, errors.New("invalid user data")
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*entities.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *UserService) GetUserByLineID(ctx context.Context, lineUserID string) (*entities.User, error) {
	return s.userRepo.GetByLineUserID(ctx, lineUserID)
}

func (s *UserService) UpdateUser(ctx context.Context, user *entities.User) error {
	if !user.IsValid() {
		return errors.New("invalid user data")
	}
	return s.userRepo.Update(ctx, user)
}
