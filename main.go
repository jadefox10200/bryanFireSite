package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	ListenPort string
}

func main() {
	config := loadConfiguration()
	
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.SetTrustedProxies(nil)

	setupRoutes(engine)

	serverAddr := fmt.Sprintf(":%s", config.ListenPort)
	log.Printf("Bryan Fire Safety server listening on %s", serverAddr)
	
	if err := engine.Run(serverAddr); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func loadConfiguration() *ServerConfig {
	listenPort := os.Getenv("PORT")
	if listenPort == "" {
		listenPort = "8080"
	}

	return &ServerConfig{
		ListenPort: listenPort,
	}
}

func setupRoutes(engine *gin.Engine) {
	engine.Use(preventDirectoryTraversal())
	
	engine.GET("/", serveHomepage)
	engine.GET("/health", healthCheck)
	
	assetHandler := createAssetHandler()
	engine.GET("/styles.css", assetHandler)
	engine.GET("/script.js", assetHandler)
	engine.GET("/:filename", assetHandler)
}

func serveHomepage(ctx *gin.Context) {
	ctx.File("./index.html")
}

func healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"healthy": true})
}

func createAssetHandler() gin.HandlerFunc {
	allowedExtensions := map[string]bool{
		".css": true, ".js": true, ".png": true, 
		".jpg": true, ".jpeg": true, ".svg": true,
		".webp": true, ".ico": true, ".gif": true,
	}
	
	return func(ctx *gin.Context) {
		requestedFile := ctx.Param("filename")
		if requestedFile == "" {
			requestedFile = ctx.Request.URL.Path[1:]
		}
		
		fileExt := path.Ext(requestedFile)
		if !allowedExtensions[fileExt] {
			ctx.Status(http.StatusNotFound)
			return
		}
		
		ctx.File("./" + requestedFile)
	}
}

func preventDirectoryTraversal() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestedPath := ctx.Request.URL.Path
		
		cleanPath := path.Clean(requestedPath)
		if strings.Contains(cleanPath, "..") || cleanPath != requestedPath {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		
		pathSegments := strings.Split(strings.Trim(requestedPath, "/"), "/")
		for _, segment := range pathSegments {
			if strings.HasPrefix(segment, ".") {
				ctx.AbortWithStatus(http.StatusForbidden)
				return
			}
		}
		
		ctx.Next()
	}
}
