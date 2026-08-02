package service

import (
	"context"
	"errors"

	"github.com/murashi19/koda-b8-backend/internal/config"
	"github.com/murashi19/koda-b8-backend/internal/dto/users"
	"github.com/murashi19/koda-b8-backend/internal/models"
	"github.com/murashi19/koda-b8-backend/internal/repo"
)

type UserService struct {
	repo             *repo.UserRepo
	refreshTokenRepo *repo.RefreshTokenRepo
	cfg              *config.Config
}

func NewUserService(
	repo *repo.UserRepo,
	refreshTokenRepo *repo.RefreshTokenRepo,
	cfg *config.Config,
) *UserService {
	return &UserService{
		repo:             repo,
		refreshTokenRepo: refreshTokenRepo,
		cfg:              cfg,
	}
}

// GET /users
func (s *UserService) GetAll(ctx context.Context) ([]*models.User, error) {
	return s.repo.FindAll(ctx)
}

// GET /users/:id
func (s *UserService) GetByID(ctx context.Context, id int64) (*models.User, error) {
	return s.repo.FindByID(ctx, id)
}

// PUT /users/:id (Admin)
func (s *UserService) Update(ctx context.Context, detail *models.UserDetail) error {
	return s.repo.Update(ctx, detail)
}

// GET /me
func (s *UserService) GetMyProfile(ctx context.Context, userID int64) (*models.UserDetail, error) {
	return s.repo.FindDetailByID(ctx, userID)
}

// PATCH /me
func (s *UserService) UpdateMyProfile(ctx context.Context, userID int64, req *users.UpdateProfileRequest) error {

	detail, err := s.repo.FindDetailByID(ctx, userID)
	if err != nil {
		return err
	}

	// Jika email berubah, cek apakah sudah dipakai user lain
	if detail.Email != req.Email {
		exists, err := s.repo.ExistsByEmail(ctx, req.Email)
		if err != nil {
			return err
		}

		if exists {
			return errors.New("email already exists")
		}

		detail.Email = req.Email
	}

	detail.Profile.FullName = req.FullName
	detail.Profile.PhoneNumber = req.PhoneNumber
	detail.Profile.BirthDate = req.BirthDate
	detail.Profile.Gender = req.Gender

	return s.repo.Update(ctx, detail)
}

// PATCH /me/avatar
func (s *UserService) UpdateAvatar(ctx context.Context, userID int64, avatar string) error {

	profile := &models.UserProfile{
		UserID: userID,
		Avatar: &avatar,
	}

	return s.repo.UpdateAvatar(ctx, profile)
}

// DELETE /users/:id
func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// GET /profile/:id
func (s *UserService) GetDetailByID(ctx context.Context, id int64) (*models.UserDetail, error) {
	return s.repo.FindDetailByID(ctx, id)
}
