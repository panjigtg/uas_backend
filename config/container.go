package config

import (
	repo "uas/app/repository"
	"database/sql"
	mgodriver "go.mongodb.org/mongo-driver/mongo"

	"uas/app/services"
)


type Container struct {
	AuthService 		*services.AuthService
	UserService 		*services.UserService
	StudentService 		*services.StudentService
	AchievementService	*services.AchievementService
}

// Dependency Injection Container
func BuildContainer(db *sql.DB, mongoDB *mgodriver.Database) *Container {

	// REPOSITORIES
	authRepo := repo.NewAuthRepo(db)
	userRepo := repo.NewUserRepo(db)
	studentRepo := repo.NewStudentRepo(db)
    lecturerRepo := repo.NewLecturerRepo(db)
	achievementRefRepo := repo.NewAchievementReferenceRepository(db)


	achievementMongoRepo := repo.NewAchievementMongoRepository(
		mongoDB.Collection("achievements"),
	)

	// SERVICES
	authService := services.NewAuthService(authRepo)
    userService := services.NewUserService(db, userRepo, studentRepo, lecturerRepo)
	studentService := services.NewStudentService(
		db,
		studentRepo,
		lecturerRepo,
		achievementRefRepo,
		achievementMongoRepo,
	)

	achievementService := services.NewAchievementService(
		studentRepo,
		achievementMongoRepo,
		achievementRefRepo,
		lecturerRepo,
		userRepo,
	)


	return &Container{
		AuthService: authService,
		UserService: userService,
		StudentService:  studentService,
		AchievementService: achievementService,
	}
}
