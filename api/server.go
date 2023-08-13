package api

import (
	db "likesApi/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
	client *redis.Client
}

// Create new instance and handle routing
func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	server.client = redisClient

	gin.SetMode(gin.ReleaseMode)
	// add routes to the router
	router.POST("/users", server.createUser)
	router.GET("/user/:id", server.getUser)
	// router.GET("/users", server.getAllUsers)
	router.POST("/likes", server.addLike)
	router.GET("/likes/:user_id/:content_id", server.getLike)
	router.GET("/likes/totallikes/:content_id", server.getTotalLikes)

	server.router = router
	return server
}

// Start runs the HTTP server on the given port
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
