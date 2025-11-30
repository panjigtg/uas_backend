package config

import (
	"database/sql"

	"uas/app/repository/psql"
	"uas/app/services"
)


type Container struct {
	AuthService *services.AuthService
	UserService *services.UserService
}

// Dependency Injection Container
func BuildContainer(db *sql.DB) *Container {

	// REPOSITORIES
	authRepo := psql.NewAuthRepo(db)
	userRepo := psql.NewUserRepo(db)

	// SERVICES
	authService := services.NewAuthService(authRepo)
	userService := services.NewUserService(userRepo)

	return &Container{
		AuthService: authService,
		UserService: userService,
	}
}
