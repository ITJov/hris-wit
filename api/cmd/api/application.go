package main

import (
	"context" // <--- DITAMBAHKAN
	"net/http"
	"time"

	"github.com/wit-id/blueprint-backend-go/common/echohttp"
	"github.com/wit-id/blueprint-backend-go/common/httpservice"
	"github.com/wit-id/blueprint-backend-go/toolkit/db/postgres"
	"github.com/wit-id/blueprint-backend-go/toolkit/log"

	"github.com/wit-id/blueprint-backend-go/common/constants"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Variabel Global
var (
	e         *echo.Echo
	s         *httpservice.Service
	appConfig *viper.Viper
)

func init() {
	log.Println("DEBUG: Vercel Function init() called. Setting up application...")

	var err error

	logger, err := log.NewFromConfig(nil, "log")
	if err != nil {
		log.Fatalf("ERROR: Failed to initialize logger in init(): %v", err)
	}
	logger.Set()
	setDefaultTimezone()

	appConfig, err = envConfigVariable("config.yaml")
	if err != nil {
		log.Fatalf("ERROR: Failed to load config in init(): %v", err)
	}
	log.Println("DEBUG: Config loaded.")

	mainDB, err := postgres.NewFromConfig(appConfig, "db")
	if err != nil {
		log.Fatalf("ERROR: Failed to setup database in init(): %v", err)
	}
	log.Println("DEBUG: Database connected.")

	s = httpservice.NewService(mainDB, appConfig)
	log.Println("DEBUG: Services initialized.")

	e = echo.New() // Membuat instance Echo di sini

	// Atur penanganan error Echo Anda di sini
	e.HTTPErrorHandler = echohttp.HandleEchoError(appConfig)
	e.Static("/uploads", "uploads") // Jika ini penting untuk Vercel

	// --- KONFIGURASI CORS (di application.go init() sebagai setup global) ---
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://hris-wit-chernojovs-projects.vercel.app"}, // Ganti dengan domain frontend Anda
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS, echo.HEAD, echo.PATCH},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", constants.DefaultAllowHeaderToken, constants.DefaultAllowHeaderRefreshToken},
		AllowCredentials: true,
	}))
	log.Println("DEBUG: CORS middleware set in init().")

	// --- Panggil fungsi untuk mendaftarkan semua routes ---
	// Ini adalah perubahan utama. Memanggil fungsi yang sudah diadaptasi.
	echohttp.SetupAllRoutes(context.Background(), s, appConfig, e) // <--- Panggilan ke fungsi baru

	log.Println("DEBUG: Echo server initialized and all routes registered.")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("DEBUG: Request received by Vercel Handler.")
	e.ServeHTTP(w, r)
}

func main() {
	log.Println("DEBUG: Running application locally via main()...")
	log.Println("DEBUG: Local setup complete, starting local server.")
	port := appConfig.GetString("restapi.port")
	if port == "" {
		port = "6969"
	}
	log.Printf("DEBUG: Local server starting on :%s", port)
	e.Start(":" + port)
}

func setDefaultTimezone() {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		loc = time.Now().Location()
	}

	time.Local = loc
}

func envConfigVariable(filePath string) (cfg *viper.Viper, err error) {
	cfg = viper.New()
	cfg.SetConfigFile(filePath)

	if err = cfg.ReadInConfig(); err != nil {
		err = errors.Wrap(err, "Error while reading config file")

		return
	}

	return
}
