package main

import (
	"log"

	"example.com/m/internal/api/v1/adapters/controllers"
	"example.com/m/internal/api/v1/adapters/repositories"
	"example.com/m/internal/api/v1/core/application/services/auth_service"
	"example.com/m/internal/api/v1/core/application/services/debt_service"
	"example.com/m/internal/api/v1/core/application/services/group_service"
	"example.com/m/internal/api/v1/core/application/services/user_service"
	"example.com/m/internal/api/v1/infrastructure/cache"
	database "example.com/m/internal/api/v1/infrastructure/database"
	"example.com/m/internal/api/v1/infrastructure/middlewares"
	"example.com/m/internal/api/v1/infrastructure/prom"
	"example.com/m/internal/api/v1/infrastructure/router"
	"example.com/m/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No all env var avaiable")
	}
}

func main() {
	loadEnv()
	config.InitConfig()
	database.ConnectToDatabase()
	cache.ConnectToRedis()
	prom.RegisterPrometheusMetrics()

	defer database.Db.Close()

	userRepository := repositories.NewUserRepository(database.Db)
	tokenRepository := repositories.NewTokenRepository(cache.Redis)
	debtRepository := repositories.NewDebtRepository(database.Db)
	groupRepository := repositories.NewGroupRepository(database.Db)
	userService := user_service.NewUserService(userRepository)
	authService := auth_service.NewAuthService(userService, tokenRepository)
	groupService := group_service.NewGroupService(groupRepository, userRepository, debtRepository)
	debtService := debt_service.NewDebtService(debtRepository, userService, groupService)
	authMiddleware := middlewares.NewAuthMiddleware(authService)
	debtController := controllers.NewDebtController(debtService)
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)
	groupController := controllers.NewGroupController(groupService)
	metricController := controllers.NewMetricController()
	engine := gin.Default()

	router.BindRoutes(engine, authMiddleware, userController, authController, metricController, debtController, groupController)

	engine.Run(":8000")
}
