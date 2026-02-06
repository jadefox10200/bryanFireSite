package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

type ServerConfig struct {
	ListenPort     string
	MailServer     string
	MailServerPort int
	MailUsername   string
	MailPassword   string
	RecipientAddr  string
}

type FormSubmission struct {
	FullName        string `form:"name" binding:"required"`
	ContactEmail    string `form:"email" binding:"required,email"`
	PhoneNumber     string `form:"phone"`
	RequestDetails  string `form:"message" binding:"required"`
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

	setupRoutes(engine, config)

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

	mailServerPort := 587
	if os.Getenv("SMTP_PORT") == "465" {
		mailServerPort = 465
	}

	return &ServerConfig{
		ListenPort:     listenPort,
		MailServer:     os.Getenv("SMTP_HOST"),
		MailServerPort: mailServerPort,
		MailUsername:   os.Getenv("SMTP_USER"),
		MailPassword:   os.Getenv("SMTP_PASS"),
		RecipientAddr:  os.Getenv("TO_EMAIL"),
	}
}

func setupRoutes(engine *gin.Engine, config *ServerConfig) {
	engine.Use(preventDirectoryTraversal())
	
	engine.GET("/", serveHomepage)
	engine.GET("/health", healthCheck)
	
	assetHandler := createAssetHandler()
	engine.GET("/styles.css", assetHandler)
	engine.GET("/script.js", assetHandler)
	engine.GET("/:filename", assetHandler)
	
	engine.POST("/contact", createContactHandler(config))
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

func createContactHandler(config *ServerConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var submission FormSubmission
		
		if err := ctx.ShouldBind(&submission); err != nil {
			log.Printf("Invalid form data received: %v", err)
			ctx.String(http.StatusBadRequest, 
				"Please complete all required fields. Call us at 408-293-1077 for immediate assistance.")
			return
		}
		
		err := dispatchEmail(config, &submission)
		if err != nil {
			log.Printf("Email delivery failed: %v", err)
			ctx.String(http.StatusInternalServerError,
				"Unable to send your message at this time. Please call 408-293-1077 or email info@bryanfire.com directly.")
			return
		}
		
		ctx.String(http.StatusOK, 
			fmt.Sprintf("Thank you %s! We received your request and will respond shortly.", submission.FullName))
	}
}

func dispatchEmail(config *ServerConfig, submission *FormSubmission) error {
	if config.MailServer == "" || config.RecipientAddr == "" {
		return fmt.Errorf("mail configuration incomplete")
	}
	
	msg := gomail.NewMessage()
	msg.SetHeader("From", config.MailUsername)
	msg.SetHeader("To", config.RecipientAddr)
	msg.SetHeader("Reply-To", submission.ContactEmail)
	msg.SetHeader("Subject", fmt.Sprintf("Service Request from %s", submission.FullName))
	
	emailBody := buildEmailContent(submission)
	msg.SetBody("text/plain", emailBody)
	
	dialer := gomail.NewDialer(
		config.MailServer,
		config.MailServerPort,
		config.MailUsername,
		config.MailPassword,
	)
	
	if config.MailServerPort == 465 {
		dialer.SSL = true
	} else {
		dialer.TLSConfig = &tls.Config{ServerName: config.MailServer}
	}
	
	if err := dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("email transmission error: %w", err)
	}
	
	log.Printf("Service request email delivered from %s <%s>", submission.FullName, submission.ContactEmail)
	return nil
}

func buildEmailContent(submission *FormSubmission) string {
	var content strings.Builder
	
	content.WriteString("=== NEW SERVICE REQUEST ===\n\n")
	content.WriteString(fmt.Sprintf("Customer Name: %s\n", submission.FullName))
	content.WriteString(fmt.Sprintf("Email Address: %s\n", submission.ContactEmail))
	
	if submission.PhoneNumber != "" {
		content.WriteString(fmt.Sprintf("Phone Number: %s\n", submission.PhoneNumber))
	}
	
	content.WriteString("\nService Request Details:\n")
	content.WriteString(strings.Repeat("-", 50))
	content.WriteString("\n")
	content.WriteString(submission.RequestDetails)
	content.WriteString("\n")
	content.WriteString(strings.Repeat("-", 50))
	content.WriteString("\n\n")
	content.WriteString("Submitted via bryanfire.com contact form\n")
	
	return content.String()
}
