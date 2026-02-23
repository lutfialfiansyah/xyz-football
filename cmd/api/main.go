package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"xyz-football-api/internal/config"
	"xyz-football-api/internal/middleware"

	"xyz-football-api/internal/modules/admin"
	"xyz-football-api/internal/pkg/storage"
	"xyz-football-api/internal/pkg/storage/local"

	adminHandler "xyz-football-api/internal/modules/admin/delivery/http"
	adminRepo "xyz-football-api/internal/modules/admin/repository/postgres"
	adminUsecase "xyz-football-api/internal/modules/admin/usecase"

	teamHandler "xyz-football-api/internal/modules/team/delivery/http"
	teamRepo "xyz-football-api/internal/modules/team/repository/postgres"
	teamUsecase "xyz-football-api/internal/modules/team/usecase"

	playerHandler "xyz-football-api/internal/modules/player/delivery/http"
	playerRepo "xyz-football-api/internal/modules/player/repository/postgres"
	playerUsecase "xyz-football-api/internal/modules/player/usecase"

	matchHandler "xyz-football-api/internal/modules/match/delivery/http"
	matchRepo "xyz-football-api/internal/modules/match/repository/postgres"
	matchUsecase "xyz-football-api/internal/modules/match/usecase"

	reportHandler "xyz-football-api/internal/modules/report/delivery/http"
	reportRepo "xyz-football-api/internal/modules/report/repository/postgres"
	reportUsecase "xyz-football-api/internal/modules/report/usecase"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Initialize Config & Database
	cfg := config.LoadConfig()

	// Seed Super Admin if not exists
	seedAdmin(cfg)

	// Dependency Injection
	// Repositories
	aRepo := adminRepo.NewAdminRepository(cfg.DB)
	tRepo := teamRepo.NewTeamRepository(cfg.DB)
	pRepo := playerRepo.NewPlayerRepository(cfg.DB)
	mRepo := matchRepo.NewMatchRepository(cfg.DB)
	rRepo := reportRepo.NewReportRepository(cfg.DB)

	// Usecases
	var storageProvider storage.Provider
	if cfg.UploadStorage == "local" {
		storageProvider = local.NewLocalStorage(cfg.AppURL)
	} else {
		// Fallback or other providers
		storageProvider = local.NewLocalStorage(cfg.AppURL)
	}

	aUseCase := adminUsecase.NewAdminUsecase(aRepo, cfg.JWTKey, cfg.JWTExpiration)
	tUseCase := teamUsecase.NewTeamUsecase(tRepo, storageProvider)
	pUseCase := playerUsecase.NewPlayerUsecase(pRepo, tRepo)
	mUseCase := matchUsecase.NewMatchUsecase(mRepo, tRepo, pRepo)
	rUseCase := reportUsecase.NewReportUsecase(rRepo)

	// Handlers
	aHandler := adminHandler.NewAdminHandler(aUseCase)
	tHandler := teamHandler.NewTeamHandler(tUseCase)
	pHandler := playerHandler.NewPlayerHandler(pUseCase)
	mHandler := matchHandler.NewMatchHandler(mUseCase)
	rHandler := reportHandler.NewReportHandler(rUseCase)

	// Router Setup
	r := gin.Default()

	// Static Routes
	r.Static("/uploads", "./uploads")

	// Public Routes
	v1 := r.Group("/api/v1")
	{
		v1.POST("/login", aHandler.Login)
	}

	// Protected Routes
	protected := v1.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWTKey))
	{
		// Validation Super Admin route
		super := protected.Group("/")
		super.Use(middleware.RoleMiddleware("admin"))
		{
			// Example Route for super admin (if needed)
		}

		// Team Routes
		teams := protected.Group("/teams")
		{
			teams.POST("/", tHandler.Create)
			teams.GET("/", tHandler.GetAll)
			teams.GET("/:id", tHandler.GetByID)
			teams.PUT("/:id", tHandler.Update)
			teams.DELETE("/:id", tHandler.Delete)
		}

		// Player Routes
		players := protected.Group("/players")
		{
			players.POST("/", pHandler.Create)
			players.GET("/", pHandler.GetAll)
			players.GET("/:id", pHandler.GetByID)
			players.PUT("/:id", pHandler.Update)
			players.DELETE("/:id", pHandler.Delete)
		}

		// Match Routes
		matches := protected.Group("/matches")
		{
			matches.POST("/", mHandler.CreateMatch)
			matches.GET("/", mHandler.GetAllMatches)
			matches.GET("/:id", mHandler.GetMatchByID)
			matches.PUT("/:id/status", mHandler.ChangeMatchStatus)
			matches.PUT("/:id/score", mHandler.ReportMatchScore)
			matches.DELETE("/:id", mHandler.DeleteMatch)

			// Match Events
			matches.POST("/:id/events", mHandler.AddMatchEvent)
			matches.GET("/:id/events", mHandler.GetMatchEvents)
		}

		// Reporting Routes
		reports := protected.Group("/reports")
		{
			reports.GET("/matches", rHandler.GetMatchReports)
		}
	}

	// Server Setup
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Graceful Shutdown
	go func() {
		log.Printf("Starting server on port %s...", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func seedAdmin(cfg *config.Config) {
	var adminUser admin.Admin
	if err := cfg.DB.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		cfg.DB.Create(&admin.Admin{
			Username:     "admin",
			PasswordHash: string(hashedPass),
			Role:         "admin",
		})
		log.Println("Seeded initial admin user (username: admin, password: password)")
	}
}
