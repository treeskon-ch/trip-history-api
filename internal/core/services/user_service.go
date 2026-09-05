package services

import (
	"context"

	"trip-history-api/internal/core/domain"
	"trip-history-api/internal/core/ports"
)

type userService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) ports.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) RegisterUser(ctx context.Context, email, password, name, role string) (domain.User, string, error) {
	if role == "" {
		role = "driver"
	}
	
	// 1. Create user in Firebase Auth
	uid, err := s.userRepo.CreateUserInAuth(ctx, email, password, name)
	if err != nil {
		return domain.User{}, "", err
	}

	// 2. Save user profile to Firestore
	user := domain.User{
		UserID: uid,
		Email:  email,
		Name:   name,
		Role:   role,
	}
	err = s.userRepo.SaveUser(ctx, user)
	if err != nil {
		return domain.User{}, "", err
	}
	
	// 3. Generate Email Verification Link
	link, err := s.userRepo.GenerateEmailVerificationLink(ctx, email)
	if err != nil {
		// Even if link generation fails, user is created. We can just return empty link or error.
		return user, "", nil 
	}

	return user, link, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return s.userRepo.GetAllUsers(ctx)
}

func (s *userService) DeleteUser(ctx context.Context, uid string) error {
	err := s.userRepo.DeleteUserInFirestore(ctx, uid)
	if err != nil {
		return err
	}
	return s.userRepo.DeleteUserInAuth(ctx, uid)
}

func (s *userService) UpdateFCMToken(ctx context.Context, uid string, token string) error {
	return s.userRepo.UpdateFCMToken(ctx, uid, token)
}
