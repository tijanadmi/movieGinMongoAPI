package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/tijanadmi/movieginmongoapi/db/mock"
	"github.com/tijanadmi/movieginmongoapi/models"
	"github.com/tijanadmi/movieginmongoapi/repository"
	"github.com/tijanadmi/movieginmongoapi/token"
	"github.com/tijanadmi/movieginmongoapi/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type movieMatcher struct {
	expected *models.Movie
}

func (m movieMatcher) Matches(x interface{}) bool {
	actual, ok := x.(*models.Movie)
	if !ok {
		return false
	}

	return actual.Title == m.expected.Title &&
		actual.Duration == m.expected.Duration &&
		actual.Genre == m.expected.Genre &&
		actual.Directors == m.expected.Directors &&
		actual.Actors == m.expected.Actors &&
		actual.Screening.Truncate(time.Second).Equal(m.expected.Screening.Truncate(time.Second)) &&
		actual.Plot == m.expected.Plot &&
		actual.Poster == m.expected.Poster
}

func (m movieMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expected)
}
func TestSearchMovieAPI(t *testing.T) {
	username := util.RandomOwner()
	role := util.UserRole
	movie := randomMovie()
	//movies := []models.Movie{movie}

	testCases := []struct {
		name          string
		movieID       primitive.ObjectID
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			movieID: movie.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMovie(gomock.Any(), gomock.Eq(movie.ID.Hex())).
					Times(1).
					Return(&movie, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchMovie(t, recorder.Body, movie)
			},
		},
		{
			name:    "NoAuthorization",
			movieID: movie.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMovie(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:    "NotFound",
			movieID: movie.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, role, time.Minute)
			},

			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMovie(gomock.Any(), gomock.Eq(movie.ID.Hex())).
					Times(1).
					Return(nil, repository.ErrMovieNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:    "InternalError",
			movieID: movie.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetMovie(gomock.Any(), gomock.Eq(movie.ID.Hex())).
					Times(1).
					Return(nil, mongo.ErrClientDisconnected)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/movies/%s", tc.movieID.Hex())
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func TestListMoviesAPI(t *testing.T) {
	username := util.RandomOwner()
	role := util.UserRole
	n := 5
	movies := make([]models.Movie, n)
	for i := 0; i < n; i++ {
		movies[i] = randomMovie()
	}

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					SearchMovies(gomock.Any(), gomock.Eq("0")).
					Times(1).
					Return(movies, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				//requireBodyMatchHalls(t, recorder.Body, halls)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					SearchMovies(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					SearchMovies(gomock.Any(), gomock.Eq("0")).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/movies"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestInsertMovieAPI(t *testing.T) {
	username := util.RandomOwner()
	role := util.UserRole
	movie := randomMovie()

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"title":     movie.Title,
				"duration":  movie.Duration,
				"genre":     movie.Genre,
				"directors": movie.Directors,
				"actors":    movie.Actors,
				"screening": movie.Screening,
				"plot":      movie.Plot,
				"poster":    movie.Poster,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := &models.Movie{
					Title:     movie.Title,
					Duration:  movie.Duration,
					Genre:     movie.Genre,
					Directors: movie.Directors,
					Actors:    movie.Actors,
					Screening: movie.Screening,
					Plot:      movie.Plot,
					Poster:    movie.Poster,
				}
				store.EXPECT().AddMovie(gomock.Any(), movieMatcher{expected: arg}).Times(1).Return(&movie, nil)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchMovie(t, recorder.Body, movie)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"title":     movie.Title,
				"duration":  movie.Duration,
				"genre":     movie.Genre,
				"directors": movie.Directors,
				"actors":    movie.Actors,
				"screening": movie.Screening,
				"plot":      movie.Plot,
				"poster":    movie.Poster,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					AddMovie(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidRows",
			body: gin.H{
				"title":     movie.Title,
				"duration":  "invalid",
				"genre":     movie.Genre,
				"directors": movie.Directors,
				"actors":    movie.Actors,
				"screening": movie.Screening,
				"plot":      movie.Plot,
				"poster":    movie.Poster,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					AddMovie(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/movies"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateMovieAPI(t *testing.T) {
	username := util.RandomOwner()
	role := util.UserRole
	movie := randomMovie()

	testCases := []struct {
		name          string
		movieID       string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			movieID: movie.ID.Hex(),
			body: gin.H{
				"id":        movie.ID.Hex(),
				"title":     movie.Title,
				"duration":  movie.Duration,
				"genre":     movie.Genre,
				"directors": movie.Directors,
				"actors":    movie.Actors,
				"screening": movie.Screening,
				"plot":      movie.Plot,
				"poster":    movie.Poster,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := &models.Movie{
					ID:        movie.ID,
					Title:     movie.Title,
					Duration:  movie.Duration,
					Genre:     movie.Genre,
					Directors: movie.Directors,
					Actors:    movie.Actors,
					Screening: movie.Screening,
					Plot:      movie.Plot,
					Poster:    movie.Poster,
				}
				store.EXPECT().UpdateMovie(gomock.Any(), gomock.Eq(movie.ID.Hex()), gomock.Any()).Times(1).Return(arg, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				// requireBodyMatchMovie(t, recorder.Body, movie)
			},
		},
		{
			name:    "NoAuthorization",
			movieID: movie.ID.Hex(),
			body: gin.H{
				"id":        movie.ID.Hex(),
				"title":     movie.Title,
				"duration":  movie.Duration,
				"genre":     movie.Genre,
				"directors": movie.Directors,
				"actors":    movie.Actors,
				"screening": movie.Screening,
				"plot":      movie.Plot,
				"poster":    movie.Poster,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// Do not set up authorization
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateMovie(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:    "InternalError",
			movieID: movie.ID.Hex(),
			body: gin.H{
				"id":        movie.ID.Hex(),
				"title":     movie.Title,
				"duration":  movie.Duration,
				"genre":     movie.Genre,
				"directors": movie.Directors,
				"actors":    movie.Actors,
				"screening": movie.Screening,
				"plot":      movie.Plot,
				"poster":    movie.Poster,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := &models.Movie{
					ID:        movie.ID,
					Title:     movie.Title,
					Duration:  movie.Duration,
					Genre:     movie.Genre,
					Directors: movie.Directors,
					Actors:    movie.Actors,
					Screening: movie.Screening,
					Plot:      movie.Plot,
					Poster:    movie.Poster,
				}
				store.EXPECT().UpdateMovie(gomock.Any(), gomock.Eq(movie.ID.Hex()), gomock.Any()).Times(1).Return(arg, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			// url := "/movies/" + tc.movieID
			url := fmt.Sprintf("/movies/%s", tc.movieID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

// func TestUpdateMovieAPI(t *testing.T) {
// 	username := util.RandomOwner()
// 	role := util.UserRole
// 	movie := randomMovie()

// 	testCases := []struct {
// 		name          string
// 		movieID       string
// 		body          gin.H
// 		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name:    "OK",
// 			movieID: movie.ID.Hex(),
// 			body: gin.H{
// 				"id":        movie.ID.Hex(),
// 				"title":     movie.Title,
// 				"duration":  movie.Duration,
// 				"genre":     movie.Genre,
// 				"directors": movie.Directors,
// 				"actors":    movie.Actors,
// 				"screening": movie.Screening,
// 				"plot":      movie.Plot,
// 				"poster":    movie.Poster,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, role, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				arg := &models.Movie{
// 					ID:        movie.ID,
// 					Title:     movie.Title,
// 					Duration:  movie.Duration,
// 					Genre:     movie.Genre,
// 					Directors: movie.Directors,
// 					Actors:    movie.Actors,
// 					Screening: movie.Screening,
// 					Plot:      movie.Plot,
// 					Poster:    movie.Poster,
// 				}
// 				store.EXPECT().UpdateMovie(gomock.Any(), gomock.Eq(movie.ID.Hex()), movieMatcher{expected: arg}).Times(1).Return(movie, nil)

// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				//requireBodyMatchMovie(t, recorder.Body, movie)
// 			},
// 		},
// 		{
// 			name:    "NoAuthorization",
// 			movieID: movie.ID.Hex(),
// 			body: gin.H{
// 				"id":        movie.ID.Hex(),
// 				"title":     movie.Title,
// 				"duration":  movie.Duration,
// 				"genre":     movie.Genre,
// 				"directors": movie.Directors,
// 				"actors":    movie.Actors,
// 				"screening": movie.Screening,
// 				"plot":      movie.Plot,
// 				"poster":    movie.Poster,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				// Do not set up authorization
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					UpdateMovie(gomock.Any(), gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name:    "InternalError",
// 			movieID: movie.ID.Hex(),
// 			body: gin.H{
// 				"id":        movie.ID.Hex(),
// 				"title":     movie.Title,
// 				"duration":  movie.Duration,
// 				"genre":     movie.Genre,
// 				"directors": movie.Directors,
// 				"actors":    movie.Actors,
// 				"screening": movie.Screening,
// 				"plot":      movie.Plot,
// 				"poster":    movie.Poster,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, role, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				arg := models.Movie{
// 					ID:        movie.ID,
// 					Title:     movie.Title,
// 					Duration:  movie.Duration,
// 					Genre:     movie.Genre,
// 					Directors: movie.Directors,
// 					Actors:    movie.Actors,
// 					Screening: movie.Screening,
// 					Plot:      movie.Plot,
// 					Poster:    movie.Poster,
// 				}
// 				store.EXPECT().
// 					UpdateMovie(gomock.Any(), gomock.Eq(movie.ID.Hex()), gomock.Eq(arg)).
// 					Times(1).
// 					Return(models.Movie{}, sql.ErrConnDone)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			},
// 		},
// 		{
// 			name:    "InvalidRows",
// 			movieID: movie.ID.Hex(),
// 			body: gin.H{
// 				"id":        movie.ID.Hex(),
// 				"title":     movie.Title,
// 				"duration":  movie.Duration,
// 				"genre":     movie.Genre,
// 				"directors": movie.Directors,
// 				"actors":    movie.Actors,
// 				"screening": movie.Screening,
// 				"plot":      movie.Plot,
// 				"poster":    movie.Poster,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
// 				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, role, time.Minute)
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					UpdateMovie(gomock.Any(), gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := newTestServer(t, store)
// 			recorder := httptest.NewRecorder()

// 			// Marshal body data to JSON
// 			data, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			//url := "/movies/" + tc.movieID
// 			url := fmt.Sprintf("/movies/%s", tc.movieID)
// 			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
// 			require.NoError(t, err)

// 			tc.setupAuth(t, request, server.tokenMaker)
// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(recorder)
// 		})
// 	}
// }

func TestDeleteMovieAPI(t *testing.T) {
	username := util.RandomOwner()
	role := util.UserRole
	movie := randomMovie()
	movieID := movie.ID.Hex()

	testCases := []struct {
		name          string
		movieID       string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			movieID: movieID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteMovie(gomock.Any(), gomock.Eq(movieID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchResponse(t, recorder.Body, apiResponse{Message: "Movie has been deleted"})
			},
		},
		{
			name:    "InternalError",
			movieID: movieID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteMovie(gomock.Any(), gomock.Eq(movieID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				requireBodyMatchErrorResponse(t, recorder.Body, apiErrorResponse{Error: sql.ErrConnDone.Error()})
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/movies/" + tc.movieID
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

// func randomMovie() models.Movie {
// 	objectID := primitive.NewObjectID()
// 	return models.Movie{
// 		ID:        objectID,
// 		Title:     "Titanik",
// 		Duration:  int32(util.RandomInt(100, 250)),
// 		Genre:     util.RandomString(5),
// 		Directors: util.RandomString(10),
// 		Actors:    util.RandomString(20),
// 		Screening: time.Now(),
// 		Plot:      util.RandomString(20),
// 		Poster:    util.RandomString(22),
// 	}
// }

func randomMovie() models.Movie {
	objectID := primitive.NewObjectID()
	return models.Movie{
		ID:        objectID,
		Title:     "Titanik",
		Duration:  132,
		Genre:     "drama",
		Directors: "Kameron",
		Actors:    "Leonardo Di Kaprio, Kejt Vinslet",
		Screening: time.Now(),
		Plot:      "nekada davno",
		Poster:    "skinuti sa interneta",
	}
}

func requireBodyMatchMovie(t *testing.T, body *bytes.Buffer, movie models.Movie) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotMovie models.Movie
	err = json.Unmarshal(data, &gotMovie)
	require.NoError(t, err)

	// Normalize the Screening time to zero out the nanosecond and location
	gotMovie.Screening = gotMovie.Screening.Truncate(time.Second)
	movie.Screening = movie.Screening.Truncate(time.Second)

	require.Equal(t, movie, gotMovie)
}

func requireBodyMatchMovies(t *testing.T, body *bytes.Buffer, movies []models.Movie) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotMovies []models.Movie
	err = json.Unmarshal(data, &gotMovies)
	require.NoError(t, err)
	require.Equal(t, movies, gotMovies)
}
