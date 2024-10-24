package cmd

import (
	"ewallet-ums/helpers"
	"ewallet-ums/internal/api"
	"ewallet-ums/internal/interfaces"
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

	userV1WithAuth := userV1.Use()
	userV1WithAuth.DELETE("/logout", dependency.MiddlewareValidateAuth, dependency.LogoutAPI.Logout)
	userV1WithAuth.PUT("/refresh-token", dependency.MiddlewareRefreshToken, dependency.RefreshTokenAPI.RefreshToken)

	err := r.Run(":" + helpers.GetEnv("PORT", ""))
	if err != nil {
		log.Fatal(err)
	}
}

type Dependency struct {
	UserRepository interfaces.IUserRepository

	HealthcheckAPI  interfaces.IHealthcheckHandler
	RegisterAPI     interfaces.IRegisterHandler
	LoginAPI        interfaces.ILoginHandler
	LogoutAPI       interfaces.ILogoutHandler
	RefreshTokenAPI interfaces.IRefreshTokenHandler

	TokenValidationAPI *api.TokenValidationHandler
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

	logoutSvc := &services.LogoutService{
		UserRepo: userRepo,
	}
	logoutAPI := &api.LogoutHandler{

		LogoutService: logoutSvc,
	}
	refreshTokenSvc := &services.RefreshTokenService{
		UserRepo: userRepo,
	}
	refreshTokenAPI := &api.RefreshTokenHandler{
		RefreshTokenService: refreshTokenSvc,
	}

	tokenValidationSvc := &services.TokenValidationService{
		UserRepo: userRepo,
	}
	tokenValidationAPI := &api.TokenValidationHandler{
		TokenValidationService: tokenValidationSvc,
	}

	return Dependency{
		UserRepository:     userRepo,
		HealthcheckAPI:     healthcheckAPI,
		RegisterAPI:        registerAPI,
		LoginAPI:           loginAPI,
		LogoutAPI:          logoutAPI,
		RefreshTokenAPI:    refreshTokenAPI,
		TokenValidationAPI: tokenValidationAPI,
	}
}
