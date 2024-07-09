package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/movieginmongoapi/repository"
	"github.com/tijanadmi/movieginmongoapi/util"
)

// GetAllReservationsForUser godoc
// @Security bearerAuth
// @Summary Get all the existing reservations for user
// @Description Get all the existing reservations for user
// @ID GetAllReservationsForUser
// @Accept  json
// @Produce  json
// @Param  username query string true "Username"
// @Success 200 {array} models.Reservation
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /reservationforuser [get]
func (server *Server) GetAllReservationsForUser(ctx *gin.Context) {

	username := ctx.Query("username")

	repertoires, err := server.store.GetAllReservationsForUser(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, repertoires)
}

// AddReservation godoc
// @Security bearerAuth
// @Summary Insert new reservation
// @Description Insert new reservation
// @ID AddReservation
// @Accept  json
// @Produce  json
// @Param reservation body models.Reservation true "Create reservation"
// @Success 201 {array} models.Reservation
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /reservation [post]
func (server *Server) AddReservation(ctx *gin.Context) {
	var req reservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apiErrorResponse{Error: "Invalid input"})
		return
	}

	dateValue, err := util.ParseDate(req.Date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, apiErrorResponse{Error: "Error parsing date"})
		return
	}

	req1 := repository.AddReservationParams{
		Username:    req.Username,
		MovieID:     req.MovieID,
		Date:        dateValue,
		Time:        req.Time,
		Hall:        req.Hall,
		ReservSeats: req.ReservSeats,
	}
	err = server.store.AddReservation(ctx, req1)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, apiResponse{Message: "Reservation added successfully"})
}

// CancelReservation godoc
// @Security bearerAuth
// @Summary Delete a single reservation
// @Description Delete a single reservation
// @ID CancelReservation
// @Accept  json
// @Produce  json
// @Param  id path string true "reservation ID"
// @Success 200 {array} apiResponse
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /reservation/{id} [delete]
func (server *Server) CancelReservation(ctx *gin.Context) {
	id := ctx.Param("id")

	err := server.store.CancelReservation(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, apiResponse{Message: "Reservation canceled successfully"})

}
