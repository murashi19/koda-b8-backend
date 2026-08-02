package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/murashi19/koda-b8-backend/internal/models"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM users
		WHERE email = $1
	)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {

	query := `
	SELECT
		u.id,
		u.email,
		u.password,
		u.role,
		up.full_name,
		u.is_verified,
		u.is_active,
		u.created_at,
		u.updated_at
	FROM users u
	JOIN user_profiles up ON up.user_id = u.id
	WHERE u.email = $1
	`

	return oneRow[models.User](ctx, r.db, query, email)
}

func (r *UserRepo) FindByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
	SELECT
		id,
		email,
		password,
		role,
		is_verified,
		is_active,
		created_at,
		updated_at
	FROM users
	WHERE id = $1
	`
	return oneRow[models.User](ctx, r.db, query, id)
}

func (r *UserRepo) FindAll(ctx context.Context) ([]*models.User, error) {
	query := `
	SELECT
		id,
		email,
		role,
		is_verified,
		is_active,
		created_at,
		updated_at
	FROM users
	ORDER BY created_at DESC
	`

	return rows[models.User](ctx, r.db, query)
}

func (r *UserRepo) Create(ctx context.Context, user *models.User, profile *models.UserProfile) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	queryUser := `
	INSERT INTO users (
		email,
		password,
		role,
		is_verified,
		is_active
	)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`
	// fmt.Println(queryUser)
	err = tx.QueryRow(ctx, queryUser, user.Email, user.Password, user.Role, user.IsVerified, user.IsActive).Scan(&user.ID)
	if err != nil {
		return err
	}

	queryProfile := `
	INSERT INTO user_profiles (
		user_id,
		full_name
	)
	VALUES ($1, $2)
	`
	_, err = tx.Exec(ctx,
		queryProfile,
		user.ID,
		profile.FullName,
	)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *UserRepo) Update(ctx context.Context, detail *models.UserDetail) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queryUser := `
	UPDATE users
	SET
		email = $1,
		role = $2,
		is_verified = $3,
		is_active = $4,
		updated_at = NOW()
	WHERE id = $5
	`

	_, err = tx.Exec(
		ctx,
		queryUser,
		detail.Email,
		detail.Role,
		detail.IsVerified,
		detail.IsActive,
		detail.ID,
	)
	if err != nil {
		return err
	}

	queryProfile := `
	UPDATE user_profiles
	SET
		full_name = $1,
		phone_number = $2,
		birth_date = $3,
		gender = $4,
		updated_at = NOW()
	WHERE user_id = $5
	`

	_, err = tx.Exec(
		ctx,
		queryProfile,
		detail.Profile.FullName,
		detail.Profile.PhoneNumber,
		detail.Profile.BirthDate,
		detail.Profile.Gender,
		detail.ID,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *UserRepo) UpdateProfile(ctx context.Context, profile *models.UserProfile) error {
	query := `
	UPDATE user_profiles
	SET
		full_name = $1,
		phone_number = $2,
		birth_date = $3,
		gender = $4,
		updated_at = NOW()
	WHERE user_id = $5
	`

	_, err := r.db.Exec(
		ctx,
		query,
		profile.FullName,
		profile.PhoneNumber,
		profile.BirthDate,
		profile.Gender,
		profile.UserID,
	)

	return err
}
func (r *UserRepo) UpdateAvatar(ctx context.Context, profile *models.UserProfile) error {
	query := `
	UPDATE user_profiles
	SET
		avatar = $1,
	WHERE user_id = $2
	`

	_, err := r.db.Exec(
		ctx,
		query,
		profile.Avatar,
		profile.UserID,
	)

	return err
}

func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		DELETE FROM user_profiles
		WHERE user_id = $1
	`, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		DELETE FROM users
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *UserRepo) FindDetailByID(ctx context.Context, id int64) (*models.UserDetail, error) {
	query := `
	SELECT
		u.id,
		u.email,
		u.password,
		u.role,
		u.is_verified,
		u.is_active,
		u.created_at,
		u.updated_at,

		up.user_id,
		up.full_name,
		up.phone_number,
		up.avatar,
		up.birth_date,
		up.gender,
		up.created_at,
		up.updated_at
	FROM users u
	INNER JOIN user_profiles up
		ON up.user_id = u.id
	WHERE u.id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	detail := &models.UserDetail{
		Profile: &models.UserProfile{},
	}

	err := row.Scan(
		&detail.ID,
		&detail.Email,
		&detail.Password,
		&detail.Role,
		&detail.IsVerified,
		&detail.IsActive,
		&detail.CreatedAt,
		&detail.UpdatedAt,

		&detail.Profile.UserID,
		&detail.Profile.FullName,
		&detail.Profile.PhoneNumber,
		&detail.Profile.Avatar,
		&detail.Profile.BirthDate,
		&detail.Profile.Gender,
		&detail.Profile.CreatedAt,
		&detail.Profile.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return detail, nil
}
