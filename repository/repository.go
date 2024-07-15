package repository

import (
	"context"
	"time"

	"github.com/tijanadmi/movieginmongoapi/models"
)

type Store interface {
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	InsertUser(ctx context.Context, user *models.User) (*models.User, error)

	InsertHall(ctx context.Context, hall *models.Hall) (*models.Hall, error)
	ListHalls(ctx context.Context) ([]models.Hall, error)
	GetHall(ctx context.Context, name string) ([]models.Hall, error)
	GetHallById(ctx context.Context, id string) (*models.Hall, error)
	UpdateHall(ctx context.Context, id string, hall models.Hall) (models.Hall, error)
	DeleteHall(ctx context.Context, id string) error

	AddMovie(ctx context.Context, movie *models.Movie) (*models.Movie, error)
	ListMovies(ctx context.Context) ([]models.Movie, error)
	GetMovie(ctx context.Context, id string) (*models.Movie, error)
	UpdateMovie(ctx context.Context, id string, movie *models.Movie) (*models.Movie, error)
	DeleteMovie(ctx context.Context, id string) error
	SearchMovies(ctx context.Context, movieId string) ([]models.Movie, error)

	AddRepertoire(ctx context.Context, repertoire *models.Repertoire) (*models.Repertoire, error)
	ListRepertoires(ctx context.Context) ([]models.Repertoire, error)
	GetRepertoire(ctx context.Context, id string) (*models.Repertoire, error)
	GetRepertoireByMovieDateTimeHall(ctx context.Context, movieId string, dateValue time.Time, timeValue string, hallValue string) (models.Repertoire, error)
	GetAllRepertoireForMovie(ctx context.Context, movieId string, startDate time.Time, endDate time.Time) ([]models.Repertoire, error)
	UpdateRepertoire(ctx context.Context, id string, repertoire models.Repertoire) (*models.Repertoire, error)
	DeleteRepertoire(ctx context.Context, id string) error
	DeleteRepertoireForMovie(ctx context.Context, movieId string) error

	InsertReservation(ctx context.Context, reservation *models.Reservation) (*models.Reservation, error)
	GetReservationById(ctx context.Context, id string) (*models.Reservation, error)
	GetAllReservationsForUser(ctx context.Context, username string) ([]models.Reservation, error)
	DeleteReservation(ctx context.Context, id string) error

	AddReservation(ctx context.Context, req AddReservationParams) (*models.Reservation, error)
	CancelReservation(ctx context.Context, resId string) error
}
