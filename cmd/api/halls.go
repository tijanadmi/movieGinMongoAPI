package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/moveginmongo/models"
)

// listHalls godoc
// @Security bearerAuth
// @Summary List existing halls
// @Description Get all the existing halls
// @ID listHalls
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Hall
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /halls [get]
func (server *Server) listHalls(ctx *gin.Context) {

	halls, err := server.store.Hall.ListHalls(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, halls)
}

// searchHall godoc
// @Security bearerAuth
// @Summary List existing halls
// @Description Get the existing halls
// @ID searchHall
// @Accept  json
// @Produce  json
// @Param  name path string true "Hall name"
// @Success 200 {array} models.Hall
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /halls/{name} [get]
func (server *Server) searchHall(ctx *gin.Context) {

	name := ctx.Param("name")
	halls, err := server.store.Hall.GetHall(ctx, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, halls)
}

// InsertHall godoc
// @Security bearerAuth
// @Summary Insert new hall
// @Description Insert new hall
// @ID InsertHall
// @Accept  json
// @Produce  json
// @Param hall body models.Hall true "Create hall"
// @Success 201 {array} models.Hall
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /halls [post]
func (server *Server) InsertHall(ctx *gin.Context) {
	var hall *models.Hall
	if err := ctx.ShouldBindJSON(&hall); err != nil {
		ctx.JSON(http.StatusBadRequest, apiErrorResponse{Error: err.Error()})
		return
	}

	hall, err := server.store.Hall.InsertHall(ctx, hall)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, hall)
}

// UpdateHall godoc
// @Security bearerAuth
// @Summary Update a single hall
// @Description Update a single hall
// @ID UpdateHall
// @Accept  json
// @Produce  json
// @Param  id path string true "Hall ID"
// @Param hall body models.Hall true "Update hall"
// @Success 200 {array} models.Hall
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /halls/{id} [put]
func (server *Server) UpdateHall(ctx *gin.Context) {
	id := ctx.Param("id")
	var hall *models.Hall
	if err := ctx.ShouldBindJSON(&hall); err != nil {
		ctx.JSON(http.StatusBadRequest, apiErrorResponse{Error: err.Error()})
		return
	}

	modifiedHall, err := server.store.Hall.UpdateHall(ctx, id, *hall)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, modifiedHall)
}

// DeleteHall godoc
// @Security bearerAuth
// @Summary Delete a single hall
// @Description Delete a single hall
// @ID DeleteHall
// @Accept  json
// @Produce  json
// @Param  id path string true "Hall ID"
// @Success 200 {array} apiResponse
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /halls/{id} [delete]
func (server *Server) DeleteHall(ctx *gin.Context) {
	id := ctx.Param("id")

	err := server.store.Hall.DeleteHall(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, apiResponse{Message: "Hall has been deleted"})
}
