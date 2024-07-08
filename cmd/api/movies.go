package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/moveginmongo/models"
)

// searchMovies godoc
// @Security bearerAuth
// @Summary List existing movie by id
// @Description Get the existing movie by id
// @ID searchMovies
// @Accept  json
// @Produce  json
// @Param  id path string true "Movie ID"
// @Success 200 {array} models.Movie
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /movies/{id} [get]
func (server *Server) searchMovies(ctx *gin.Context) {

	movieId := ctx.Param("id")
	movies, err := server.store.Movie.SearchMovies(ctx, movieId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, movies)
}

// listMovies godoc
// @Security bearerAuth
// @Summary List existing movies
// @Description Get all the existing movies
// @ID listMovies
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Movie
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /movies [get]
func (server *Server) listMovies(ctx *gin.Context) {

	movies, err := server.store.Movie.SearchMovies(ctx, "0")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, movies)
}

// InsertMovie godoc
// @Security bearerAuth
// @Summary Insert new movie
// @Description Insert new movie
// @ID InsertMovie
// @Accept  json
// @Produce  json
// @Param movie body models.Movie true "Create movie"
// @Success 201 {array} models.Movie
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /movies [post]
func (server *Server) InsertMovie(ctx *gin.Context) {
	var movie *models.Movie
	if err := ctx.ShouldBindJSON(&movie); err != nil {
		ctx.JSON(http.StatusBadRequest, apiErrorResponse{Error: fmt.Sprintf(" invalid input: %s", err.Error())})
		return
	}

	movie, err := server.store.Movie.AddMovie(ctx, movie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, movie)
}

// UpdateMovie godoc
// @Security bearerAuth
// @Summary Update a single movie
// @Description Update a single movie
// @ID UpdateMovie
// @Accept  json
// @Produce  json
// @Param  id path string true "Movie ID"
// @Param movie body models.Movie true "Update hall"
// @Success 200 {array} models.Movie
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /movies/{id} [put]
func (server *Server) UpdateMovie(ctx *gin.Context) {
	id := ctx.Param("id")
	var movie *models.Movie
	if err := ctx.ShouldBindJSON(&movie); err != nil {
		ctx.JSON(http.StatusBadRequest, apiErrorResponse{Error: fmt.Sprintf(" invalid input: %s", err.Error())})
		return
	}

	modifiedMovie, err := server.store.Movie.UpdateMovie(ctx, id, *movie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, modifiedMovie)
}

// DeleteMovie godoc
// @Security bearerAuth
// @Summary Delete a single movie
// @Description Delete a single movie
// @ID DeleteMovie
// @Accept  json
// @Produce  json
// @Param  id path string true "Movie ID"
// @Success 200 {array} apiResponse
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /movies/{id} [delete]
func (server *Server) DeleteMovie(ctx *gin.Context) {
	id := ctx.Param("id")

	err := server.store.Movie.DeleteMovie(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, apiResponse{Message: "Movie has been deleted"})
}
