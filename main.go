// main.go
package main

import (
	"SimpleWeb/api"
	_ "SimpleWeb/docs" // docs is generated by Swag CLI
	"context"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// Create a new Gin router
	router := gin.Default()

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// api routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", api.Ping)                   // removed trailing slash
		v1.GET("/user/:id", api.GetUser)            // already correct
		v1.GET("/launchProcess", api.LaunchProcess) // removed trailing slash
		v1.GET("/resetDatabase", api.ResetDatabase) // removed trailing slash and fixed case
		v1.GET("/top-devices", api.TopDevices)      // removed trailing slash and fixed case
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	quit := make(chan os.Signal, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
			// Signal the application to shut down
			quit <- os.Interrupt
		}
	}()

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server exiting")
}
