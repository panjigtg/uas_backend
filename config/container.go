package config

import (
	"database/sql"

	"uas/app/repository/psql"
	"uas/app/services"
)


type Container struct {
	AuthService *services.AuthService
	UserService *services.UserService
	StudentService *services.StudentService
}

// Dependency Injection Container
func BuildContainer(db *sql.DB) *Container {

	// REPOSITORIES
	authRepo := psql.NewAuthRepo(db)
	userRepo := psql.NewUserRepo(db)
	studentRepo := psql.NewStudentRepo(db)
    lecturerRepo := psql.NewLecturerRepo(db)

	// SERVICES
	authService := services.NewAuthService(authRepo)
    userService := services.NewUserService(db, userRepo, studentRepo, lecturerRepo)
	studentService := services.NewStudentService(db, studentRepo, lecturerRepo)


	return &Container{
		AuthService: authService,
		UserService: userService,
		StudentService:  studentService,
	}
}
