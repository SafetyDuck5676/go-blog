package main

import (
	"bytes"
	"fmt"
	"go-blog/internal/database"
	"go-blog/internal/database/models"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	database.Init()

	// Set up Gin router
	router := gin.Default()

	// Load HTML templates
	router.LoadHTMLGlob("templates/*.html")

	// Routes
	router.GET("/", getPosts)              // List all posts
	router.GET("/create", showCreateForm)  // Show create form
	router.POST("/post", createPost)       // Create a new post
	router.GET("/post/:id", getPostByID)   // Get a post by ID
	router.PUT("/post/:id", updatePost)    // Update a post
	router.DELETE("/post/:id", deletePost) // Delete a post

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if PORT env variable is not set
	}
	router.Run(":" + port)
}

// Get all posts
func getPosts(c *gin.Context) {
	var posts []models.Post
	database.DB.Find(&posts)
	c.HTML(http.StatusOK, "index.html", posts)
}

// Show the create post form
func showCreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "create.html", nil)
}

// Create a new post
func createPost(c *gin.Context) {
	// Read the raw request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	if len(bodyBytes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty request body"})
		return
	}
	log.Printf("Received request body: %s", string(bodyBytes))

	// Reset the body so it can be read again by ShouldBindJSON
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var post models.Post

	// Bind and validate the JSON payload
	if err := c.ShouldBindJSON(&post); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid JSON payload: %v", err)})
		return
	}

	// Save the post to the database
	database.DB.Create(&post)

	// Return the created post
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
