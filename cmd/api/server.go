package api

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	db "github.com/tijanadmi/moveginmongo/repository"
	"github.com/tijanadmi/moveginmongo/token"
	"github.com/tijanadmi/moveginmongo/util"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	store      *db.MongoClient
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store *db.MongoClient) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	/*if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}*/

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:3000"},
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.POST("/users/login", server.loginUser)
	router.POST("/users", server.InsertUser)

	router.POST("/tokens/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/halls", server.listHalls)
	authRoutes.GET("/halls/:name", server.searchHall)
	authRoutes.PUT("/halls/:id", server.UpdateHall)
	authRoutes.POST("/halls", server.InsertHall)
	authRoutes.DELETE("/halls/:id", server.DeleteHall)

	authRoutes.GET("/movies/:id", server.searchMovies)
	authRoutes.GET("/movies", server.listMovies)
	authRoutes.PUT("/movies/:id", server.UpdateMovie)
	authRoutes.POST("/movies", server.InsertMovie)
	authRoutes.DELETE("/movies/:id", server.DeleteMovie)

	authRoutes.GET("/repertoires/:id", server.GetRepertoire)
	authRoutes.GET("/repertoires/movie", server.GetAllRepertoireForMovie)
	authRoutes.GET("/repertoires", server.ListRepertoires)
	authRoutes.PUT("/repertoires/:id", server.UpdateRepertoire)
	authRoutes.POST("/repertoires", server.AddRepertoire)
	authRoutes.DELETE("/repertoires/:id", server.DeleteRepertoire)
	authRoutes.DELETE("/repertoires/movie", server.DeleteRepertoireForMovie)

	authRoutes.POST("/reservation", server.AddReservation)
	authRoutes.DELETE("/reservation/:id", server.CancelReservation)
	authRoutes.GET("/reservationforuser", server.GetAllReservationsForUser)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
