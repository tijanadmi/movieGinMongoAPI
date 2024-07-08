package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/moveginmongo/models"
	"github.com/tijanadmi/moveginmongo/util"
)

// GetRepertoire godoc
// @Security bearerAuth
// @Summary List existing repertoire by id
// @Description Get the existing repertoire by id
// @ID GetRepertoire
// @Accept  json
// @Produce  json
// @Param  id path string true "Repertoire ID"
// @Success 200 {array} models.Repertoire
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /repertoires/{id} [get]
func (server *Server) GetRepertoire(ctx *gin.Context) {

	id := ctx.Param("id")
	repertoire, err := server.store.Repertoire.GetRepertoire(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, repertoire)
}

// GetAllRepertoireForMovie godoc
// @Security bearerAuth
// @Summary Get all the existing repertoires for the movie between startDate and endDate
// @Description Get all the existing repertoires for the movie between startDate and endDate
// @ID GetAllRepertoireForMovie
// @Accept  json
// @Produce  json
// @Param  movie_id query string true "Movie ID"
// @Param  start_date query string true "Start Date"
// @Param  end_date query string true "End Date"
// @Success 200 {array} models.Repertoire
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /repertoires/movie [get]
func (server *Server) GetAllRepertoireForMovie(ctx *gin.Context) {

	movieId := ctx.Query("movie_id")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	// Parse the date string to time.Time
	startDateValue, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing startDate"})
		return
	}

	// Parse the date string to time.Time
	endDateValue, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing endDate"})
		return
	}

	repertoires, err := server.store.Repertoire.GetAllRepertoireForMovie(ctx, movieId, startDateValue, endDateValue)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, repertoires)
}

// ListRepertoires godoc
// @Security bearerAuth
// @Summary List existing repertoires
// @Description Get the existing repertoires
// @ID ListRepertoires
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Repertoire
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /repertoires [get]
func (server *Server) ListRepertoires(ctx *gin.Context) {

	repertoires, err := server.store.Repertoire.ListRepertoires(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, repertoires)
}

// AddRepertoire godoc
// @Security bearerAuth
// @Summary Insert new repertoire
// @Description Insert new repertoire
// @ID AddRepertoire
// @Accept  json
// @Produce  json
// @Param repertoire body models.Repertoire true "Create repertoire"
// @Success 201 {array} models.Repertoire
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /repertoires [post]
func (server *Server) AddRepertoire(ctx *gin.Context) {
	var repertoire *models.Repertoire
	if err := ctx.ShouldBindJSON(&repertoire); err != nil {
		ctx.JSON(http.StatusBadRequest, apiErrorResponse{Error: fmt.Sprintf(" invalid input: %s", err.Error())})
		return
	}

	// Convert date string to time.Time
	dateValue, err := util.ParseDate(repertoire.DateSt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apiErrorResponse{Error: "Error parsing date"})
		return
	}

	// Update the date field with the parsed time.Time value
	repertoire.Date = dateValue

	repertoire, err = server.store.Repertoire.AddRepertoire(ctx, repertoire)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, repertoire)
}

// UpdateRepertoire godoc
// @Security bearerAuth
// @Summary Update a single repertoire
// @Description Update a single repertoire
// @ID UpdateRepertoire
// @Accept  json
// @Produce  json
// @Param  id path string true "Repertoire ID"
// @Param repertoire body models.Repertoire true "Update repertoire"
// @Success 200 {array} models.Repertoire
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /repertoires/{id} [put]
func (server *Server) UpdateRepertoire(ctx *gin.Context) {
	id := ctx.Param("id")
	var repertoire *models.Repertoire
	if err := ctx.ShouldBindJSON(&repertoire); err != nil {
		ctx.JSON(http.StatusBadRequest, apiErrorResponse{Error: fmt.Sprintf(" invalid input: %s", err.Error())})
		return
	}
	repertoire, err := server.store.Repertoire.UpdateRepertoire(ctx, id, *repertoire)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, repertoire)
}

// DeleteRepertoire godoc
// @Security bearerAuth
// @Summary Delete a single repertoire
// @Description Delete a single repertoire
// @ID DeleteRepertoire
// @Accept  json
// @Produce  json
// @Param  id path string true "repertoire ID"
// @Success 200 {array} apiResponse
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /repertoires/{id} [delete]
func (server *Server) DeleteRepertoire(ctx *gin.Context) {
	id := ctx.Param("id")

	err := server.store.Repertoire.DeleteRepertoire(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, apiResponse{Message: fmt.Sprintf("repertoire has been deleted")})
}

// DeleteRepertoireForMovie godoc
// @Security bearerAuth
// @Summary Delete all repertoires for the movie
// @Description Delete all repertoires for the movie
// @ID DeleteRepertoireForMovie
// @Accept  json
// @Produce  json
// @Param  movie_id query string true "movie ID"
// @Success 200 {array} apiResponse
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /repertoires/movie [delete]
func (server *Server) DeleteRepertoireForMovie(ctx *gin.Context) {
	movieId := ctx.Query("movie_id")

	err := server.store.Repertoire.DeleteRepertoireForMovie(ctx, movieId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, apiResponse{Message: fmt.Sprintf("repertoire has been deleted")})
}
