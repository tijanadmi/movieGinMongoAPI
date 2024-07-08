package api

import (
	"errors"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/moveginmongo/models"
	"github.com/tijanadmi/moveginmongo/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
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

	repertoires, err := server.store.Reservation.GetAllReservationsForUser(ctx, username)
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

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	// Starts a session on the client
	session, err := server.store.Reservation.Col.Database().Client().StartSession()
	if err != nil {
		panic(err)
	}
	// Defers ending the session after the transaction is committed or ended
	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
		// Prvo čitanje filma
		movie, err := server.store.Movie.GetMovie(sessionCtx, req.MovieID)
		if err != nil {
			return nil, err
		}

		// Zatim čitanje korisnika
		user, err := server.store.Users.GetUserByUsername(sessionCtx, req.Username)
		if err != nil {
			return nil, err
		}

		// Čitanje repertoara
		var repertoire models.Repertoire
		repertoire, err = server.store.Repertoire.GetRepertoireByMovieDateTimeHall(sessionCtx, req.MovieID, dateValue, req.Time, req.Hall)
		if err != nil {
			return nil, err
		}

		numOfSeats := len(req.ReservSeats)

		// Provera dostupnih mesta
		if repertoire.NumOfTickets < repertoire.NumOfResTickets+numOfSeats {
			return nil, errors.New("not enough tickets available")
		}

		// Ažuriranje repertoara
		updatedSeats := append(repertoire.ReservSeats, req.ReservSeats...)
		sort.Strings(updatedSeats)
		repertoire.ReservSeats = updatedSeats
		repertoire.NumOfResTickets = repertoire.NumOfResTickets + numOfSeats

		// Ažuriranje repertoara u okviru transakcije
		_, err = server.store.Repertoire.UpdateRepertoire(sessionCtx, repertoire.ID.Hex(), repertoire)
		if err != nil {
			return nil, err
		}

		// Kreiranje rezervacije
		sort.Strings(req.ReservSeats)
		var reservation *models.Reservation
		reservation = &models.Reservation{
			Username:      user.Username,
			UserID:        user.ID,
			MovieID:       movie.ID,
			MovieTitle:    movie.Title,
			RepertoiresID: repertoire.ID,
			Date:          repertoire.Date,
			Time:          repertoire.Time,
			Hall:          repertoire.Hall,
			CreationDate:  time.Now(),
			ReservSeats:   req.ReservSeats,
		}

		// Unos rezervacije u okviru transakcije
		reservation, err = server.store.Reservation.InsertReservation(sessionCtx, reservation)
		if err != nil {
			return nil, err
		}

		return reservation, nil
	}, txnOptions)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}
	if result == nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: "Result is empty"})
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

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	// Starts a session on the client
	session, err := server.store.Reservation.Col.Database().Client().StartSession()
	if err != nil {
		panic(err)
	}
	// Defers ending the session after the transaction is committed or ended
	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
		reservation, err := server.store.Reservation.GetReservationById(ctx, id)
		if err != nil {
			return nil, err

		}
		repertoire, err := server.store.Repertoire.GetRepertoire(ctx, reservation.RepertoiresID.Hex())
		if err != nil {
			return nil, err

		}

		difSeats := difference(repertoire.ReservSeats, reservation.ReservSeats)
		repertoire.ReservSeats = difSeats
		repertoire.NumOfResTickets = len(difSeats)

		/*** Update repertoire ****/
		_, err = server.store.Repertoire.UpdateRepertoire(ctx, reservation.RepertoiresID.Hex(), *repertoire)
		if err != nil {
			return nil, err

		}

		/**** Delete reservation ****/
		err = server.store.Reservation.DeleteReservation(ctx, id)
		if err != nil {
			return nil, err
		}

		return 0, nil

	}, txnOptions)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}
	if result == nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: "Result is empty"})
		return
	}

	ctx.JSON(http.StatusOK, apiResponse{Message: "Reservation canceled successfully"})

}

// func (server *Server) CancelReservation(ctx *gin.Context) {
// 	id := ctx.Param("id")

// 	reservation, err := server.store.Reservation.GetReservationById(ctx, id)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
// 		return
// 	}
// 	repertoire, err := server.store.Repertoire.GetRepertoire(ctx, reservation.RepertoiresID.Hex())
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
// 		return
// 	}

// 	difSeats := difference(repertoire.ReservSeats, reservation.ReservSeats)
// 	repertoire.ReservSeats = difSeats
// 	repertoire.NumOfResTickets = len(difSeats)

// 	/*** Update repertoire ****/
// 	modifiedCount, err := server.store.Repertoire.UpdateRepertoire(ctx, reservation.RepertoiresID.Hex(), repertoire)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
// 		return
// 	}
// 	if modifiedCount == 0 {
// 		ctx.JSON(http.StatusNotFound, apiErrorResponse{Error: "repertoire not found for update"})
// 		return
// 	}

// 	/**** Delete reservation ****/
// 	deletedCount, err := server.store.Reservation.DeleteReservation(ctx, id)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
// 		return
// 	}

// 	if deletedCount == 0 {
// 		ctx.JSON(http.StatusNotFound, apiErrorResponse{Error: "reservation not found"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, apiResponse{Message: "Reservation canceled successfully"})

// }

func difference(slice1, slice2 []string) []string {
	m := make(map[string]bool)
	for _, item := range slice2 {
		m[item] = true
	}
	var diff []string
	for _, item := range slice1 {
		if !m[item] {
			diff = append(diff, item)
		}
	}
	return diff
}
