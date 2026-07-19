package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"sharing-vision-backend/config"
	"sharing-vision-backend/handlers"
	"sharing-vision-backend/models"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Total-Count")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// DispatchPostRoute handles the ambiguity where POST /article/:id is specified for both update and delete
func DispatchPostRoute(c *gin.Context) {
	action := c.Query("action")
	methodOverride := c.Query("_method")
	if action == "delete" || methodOverride == "DELETE" {
		handlers.DeletePost(c)
		return
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if len(bodyBytes) == 0 {
		handlers.DeletePost(c)
		return
	}

	var testMap map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &testMap); err != nil {
		handlers.DeletePost(c)
		return
	}

	if len(testMap) == 0 {
		handlers.DeletePost(c)
		return
	}

	handlers.UpdatePost(c)
}

// ArticleDispatcher dynamically parses /article/*path to call correct handlers
func ArticleDispatcher(c *gin.Context) {
	path := c.Param("path")
	parts := []string{}
	for _, p := range strings.Split(path, "/") {
		if p != "" {
			parts = append(parts, p)
		}
	}

	method := c.Request.Method

	// Case 1: /article or /article/
	if len(parts) == 0 {
		if method == http.MethodPost {
			handlers.CreatePost(c)
		} else {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		}
		return
	}

	// Case 2: /article/<id> (e.g. /article/5)
	if len(parts) == 1 {
		id := parts[0]
		// Injects param to context so handlers.GetPostByID / UpdatePost / DeletePost can read it
		c.Params = []gin.Param{{Key: "id", Value: id}}

		if method == http.MethodGet {
			handlers.GetPostByID(c)
		} else if method == http.MethodPut || method == http.MethodPatch {
			handlers.UpdatePost(c)
		} else if method == http.MethodDelete {
			handlers.DeletePost(c)
		} else if method == http.MethodPost {
			DispatchPostRoute(c)
		} else {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		}
		return
	}

	// Case 3: /article/<limit>/<offset> (e.g. /article/10/0)
	if len(parts) == 2 {
		limit := parts[0]
		offset := parts[1]
		// Injects params to context so handlers.GetPosts can read it
		c.Params = []gin.Param{
			{Key: "limit", Value: limit},
			{Key: "offset", Value: offset},
		}

		if method == http.MethodGet {
			handlers.GetPosts(c)
		} else {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		}
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
}

func main() {
	_ = godotenv.Load()

	config.InitDB()

	log.Println("Running AutoMigrate for posts table...")
	err := config.DB.AutoMigrate(&models.Post{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	log.Println("AutoMigrate completed successfully.")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(CORSMiddleware())

	// Dispatcher handles all /article, /article/, /article/:id, /article/:limit/:offset routes
	r.Any("/article", ArticleDispatcher)
	r.Any("/article/*path", ArticleDispatcher)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
