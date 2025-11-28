package config

import (
	"database/sql"

	"uas/app/repository/psql"
	"uas/app/services"
)


type Container struct {
	AuthService *services.AuthService
}


func BuildContainer(db *sql.DB) *Container {

	// REPOSITORIES
	authRepo := psql.NewAuthRepo(db)

	// SERVICES
	authService := services.NewAuthService(authRepo)

	return &Container{
		AuthService: authService,
	}
}
