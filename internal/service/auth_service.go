package service

import (
	"context"
	"errors"
	"fmt"

	authdto "github.com/murashi19/koda-b8-backend/internal/dto/auth"
	userdto "github.com/murashi19/koda-b8-backend/internal/dto/users"
	"github.com/murashi19/koda-b8-backend/internal/lib"
	"github.com/murashi19/koda-b8-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) Register(ctx context.Context, req *authdto.RegisterRequest) error {

	exists, err := s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:      req.Email,
		Password:   string(hashedPassword),
		Role:       models.RoleCustomer,
		IsVerified: false,
		IsActive:   true,
	}

	profile := &models.UserProfile{
		FullName: req.FullName,
	}

	return s.repo.Create(ctx, user, profile)
}

func (s *UserService) Login(ctx context.Context, req *authdto.LoginRequest) (*authdto.LoginResponse, error) {

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		fmt.Println("FindByEmail:", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := lib.GenerateAccessToken(s.cfg, user.ID, string(user.Role))
	// fmt.Println("JWT Secret:", token)
	if err != nil {
		fmt.Println("GenerateAccessToken:", err)
		return nil, err
	}

	return &authdto.LoginResponse{
		Token: token,
		User: userdto.UserResponse{
			ID:         user.ID,
			FullName:   user.Fullname,
			Email:      user.Email,
			Role:       string(user.Role),
			IsVerified: user.IsVerified,
		},
	}, nil
}
