package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"todo-backend/bots"
	"todo-backend/db"
	"todo-backend/middleware"
	"todo-backend/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("PORT not set, defaulting to %s", port)
	} else {
		log.Printf("PORT set to %s", port)
	}
	databasePath := os.Getenv("DATABASE_PATH")
	d, err := db.New(databasePath)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	sqlDB, _ := d.DB()
	sqlDB.SetMaxOpenConns(1)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	jwtExpHours := parseEnvInt("JWT_EXP_HOURS", 24)
	router.Setup(r, d, time.Duration(jwtExpHours)*time.Hour)
	log.Printf("Router setup complete")
	_, err = bots.InitBots(d)
	if err != nil {
		log.Printf("bot init error: %v", err)
	} else {
		log.Printf("Bots initialized successfully")
	}
	log.Printf("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func parseEnvInt(key string, dv int) int {
	v := os.Getenv(key)
	if v == "" {
		return dv
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return dv
	}
	return n
}
