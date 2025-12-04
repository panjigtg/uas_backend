package config

import (
	"database/sql"
	mgodriver "go.mongodb.org/mongo-driver/mongo"

	"uas/app/repository/psql"
	mongorepo "uas/app/repository/mongo"
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
	authRepo := psql.NewAuthRepo(db)
	userRepo := psql.NewUserRepo(db)
	studentRepo := psql.NewStudentRepo(db)
    lecturerRepo := psql.NewLecturerRepo(db)
	achievementRefRepo := psql.NewAchievementReferenceRepository(db)


	achievementMongoRepo := mongorepo.NewAchievementMongoRepository(
		mongoDB.Collection("achievements"),
	)

	// SERVICES
	authService := services.NewAuthService(authRepo)
    userService := services.NewUserService(db, userRepo, studentRepo, lecturerRepo)
	studentService := services.NewStudentService(db, studentRepo, lecturerRepo)
	achievementService := services.NewAchievementService(
		studentRepo,
		achievementMongoRepo,
		achievementRefRepo,
	)


	return &Container{
		AuthService: authService,
		UserService: userService,
		StudentService:  studentService,
		AchievementService: achievementService,
	}
}
