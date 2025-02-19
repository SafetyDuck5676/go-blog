package main

import (
	"bytes"
	"fmt"
	"go-blog/internal/database"
	"go-blog/internal/database/models"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize the database
	database.Init()

	// Set up Gin router
	router := gin.Default()

	// Load HTML templates
	router.LoadHTMLGlob("templates/*.html")

	// Routes
	router.GET("/", getPosts)             // List all posts
	router.GET("/create", showCreateForm) // Show create form
	router.POST("/post", createPost)      // Create a new post

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if PORT env variable is not set
	}
	router.Run(":" + port)
}

func isValidJSONBody(body []byte) bool {
	return len(body) > 0 && body[0] == '{'
}

// Get all posts
func getPosts(c *gin.Context) {
	var posts []models.Post
	database.DB.Find(&posts)

	// Render the index template with the list of posts
	c.HTML(http.StatusOK, "index.html", posts)
}

// Show the create post form
func showCreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "create.html", nil)
}

func init() {
	// Set logrus to output JSON-formatted logs
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func createPost(c *gin.Context) {
	// Read the raw request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.WithError(err).Error("Failed to read request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	if len(bodyBytes) == 0 {
		logrus.Warn("Empty request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty request body"})
		return
	}

	if !isValidJSONBody(bodyBytes) {
		logrus.Warn("Invalid JSON payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Reset the body so it can be read again by ShouldBindJSON
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var post models.Post

	if err := c.ShouldBindJSON(&post); err != nil {
		logrus.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid JSON payload: %v", err)})
		return
	}

	if err := database.DB.Create(&post).Error; err != nil {
		logrus.WithError(err).Error("Failed to save post")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save post"})
		return
	}

	logrus.Info("Post created successfully")
	c.JSON(http.StatusCreated, post)
}

// Get a post by ID
func getPostByID(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := database.DB.Where("id = ?", id).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

// Update a post
func updatePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := database.DB.Where("id = ?", id).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&post)
	c.JSON(http.StatusOK, post)
}

// Delete a post
func deletePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := database.DB.Where("id = ?", id).Delete(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Post deleted"})
}
