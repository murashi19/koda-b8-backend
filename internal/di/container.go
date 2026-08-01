package di

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/murashi19/koda-b8-backend/internal/config"
	"github.com/murashi19/koda-b8-backend/internal/handlers"
	"github.com/murashi19/koda-b8-backend/internal/repo"
	"github.com/murashi19/koda-b8-backend/internal/service"
)

type Container struct {
	db  *pgxpool.Pool
	cfg *config.Config

	userRepo         *repo.UserRepo
	refreshTokenRepo *repo.RefreshTokenRepo
	userService      *service.UserService
	authHandler      *handlers.AuthHandler
	userHandler      *handlers.UserHandler
}

func (c *Container) initDeps() {
	c.userRepo = repo.NewUserRepository(c.db)
	c.refreshTokenRepo = repo.NewRefreshTokenRepo(c.db)
	c.userService = service.NewUserService(c.userRepo, c.refreshTokenRepo, c.cfg)
	c.authHandler = handlers.NewAuthHandler(c.userService)
	c.userHandler = handlers.NewUserHandler(c.userService)
}

func (c *Container) AuthHandler() *handlers.AuthHandler {
	return c.authHandler
}

func (c *Container) UserHandler() *handlers.UserHandler {
	return c.userHandler
}

func NewContainer() (*Container, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env: %w", err)
	}

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	container := &Container{
		db:  db,
		cfg: cfg,
	}

	container.initDeps()

	return container, nil
}
func (c *Container) Close() {
	c.db.Close()
}
func (c *Container) Config() *config.Config {
	return c.cfg
}
