package cmd

import (
	"ewallet-ums/helpers"
	"ewallet-ums/internal/api"
	"ewallet-ums/internal/repository"
	"ewallet-ums/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func ServeHTTP() {
	dependency := dependencyInject()

	r := gin.Default()

	r.GET("/health", dependency.HealthcheckAPI.HealthcheckHandlerHTTP)

	userV1 := r.Group("/user/v1")
	userV1.POST("/register", dependency.RegisterAPI.Register)
	userV1.POST("/login", dependency.LoginAPI.Login)

	err := r.Run(":" + helpers.GetEnv("PORT", ""))
	if err != nil {
		log.Fatal(err)
	}
}

type Dependency struct {
	HealthcheckAPI *api.Healthcheck
	RegisterAPI    *api.RegisterHandler
	LoginAPI       *api.LoginHandler
}

func dependencyInject() Dependency {
	healthcheckSvc := &services.Healthcheck{}
	healthcheckAPI := &api.Healthcheck{
		HealthcheckServices: healthcheckSvc,
	}

	userRepo := &repository.UserRepository{
		DB: helpers.DB,
	}

	registerSvc := &services.RegisterService{
		UserRepo: userRepo,
	}
	registerAPI := &api.RegisterHandler{
		RegisterService: registerSvc,
	}

	loginSvc := &services.LoginService{
		UserRepo: userRepo,
	}
	loginAPI := &api.LoginHandler{
		LoginService: loginSvc,
	}

	return Dependency{
		HealthcheckAPI: healthcheckAPI,
		RegisterAPI:    registerAPI,
		LoginAPI:       loginAPI,
	}
}
