package repository

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/tijanadmi/movieginmongoapi/models"
	"github.com/tijanadmi/movieginmongoapi/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// TransferTxParams contains the input parameters of the transfer transaction
type AddReservationParams struct {
	Username    string    `json:"username" binding:"required"`
	MovieID     string    `json:"movieId" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
	Time        string    `json:"time" binding:"required"`
	Hall        string    `json:"hall" binding:"required"`
	ReservSeats []string  `json:"reservSeats" binding:"required"`
}

func (r *MongoStore) AddReservation(ctx context.Context, req AddReservationParams) (*models.Reservation, error) {

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	// Starts a session on the client
	session, err := r.db.Client().StartSession()
	if err != nil {
		panic(err)
	}
	// Defers ending the session after the transaction is committed or ended
	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
		// Prvo čitanje filma
		movie, err := r.GetMovie(sessionCtx, req.MovieID)
		if err != nil {
			return nil, err
		}

		// Zatim čitanje korisnika
		user, err := r.GetUserByUsername(sessionCtx, req.Username)
		if err != nil {
			return nil, err
		}

		// Čitanje repertoara
		var repertoire models.Repertoire
		repertoire, err = r.GetRepertoireByMovieDateTimeHall(sessionCtx, req.MovieID, req.Date, req.Time, req.Hall)
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
		_, err = r.UpdateRepertoire(sessionCtx, repertoire.ID.Hex(), repertoire)
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
		reservation, err = r.InsertReservation(sessionCtx, reservation)
		if err != nil {
			return nil, err
		}

		return reservation, nil
	}, txnOptions)
	if err != nil {
		return nil, err

	}
	if result == nil {
		return nil, errors.New("result is empty")
	}
	reservation, ok := result.(*models.Reservation)
	if !ok {
		return nil, errors.New("unexpected result type")
	}

	return reservation, nil
}

func (r *MongoStore) CancelReservation(ctx context.Context, resId string) error {

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	// Starts a session on the client
	session, err := r.db.Client().StartSession()
	if err != nil {
		panic(err)
	}
	// Defers ending the session after the transaction is committed or ended
	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
		reservation, err := r.GetReservationById(ctx, resId)
		if err != nil {
			return nil, err

		}
		repertoire, err := r.GetRepertoire(ctx, reservation.RepertoiresID.Hex())
		if err != nil {
			return nil, err

		}

		difSeats := util.Difference(repertoire.ReservSeats, reservation.ReservSeats)
		repertoire.ReservSeats = difSeats
		repertoire.NumOfResTickets = len(difSeats)

		/*** Update repertoire ****/
		_, err = r.UpdateRepertoire(ctx, reservation.RepertoiresID.Hex(), *repertoire)
		if err != nil {
			return nil, err

		}

		/**** Delete reservation ****/
		err = r.DeleteReservation(ctx, resId)
		if err != nil {
			return nil, err
		}

		return 0, nil

	}, txnOptions)

	if err != nil {
		return err

	}
	if result == nil {
		errors.New("result is empty")
	}
	return nil

}
