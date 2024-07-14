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
	"github.com/tijanadmi/movieginmongoapi/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestSearchMovieAPI(t *testing.T) {
	movie := randomMovie()

	testCases := []struct {
		name    string
		movieID primitive.ObjectID
		//setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			movieID: movie.ID,
			/*setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
			},*/
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetHallById(gomock.Any(), gomock.Eq(movie.ID.Hex())).
					Times(1).
					Return(&movie, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchMovie(t, recorder.Body, movie)
			},
		},
		{
			name:    "NotFound",
			movieID: movie.ID,
			/*setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
			},*/

			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetHallById(gomock.Any(), gomock.Eq(movie.ID.Hex())).
					Times(1).
					Return(nil, repository.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:    "InternalError",
			movieID: movie.ID,
			/*setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
			},*/
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetHallById(gomock.Any(), gomock.Eq(movie.ID.Hex())).
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

			//tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func TestListMoviesAPI(t *testing.T) {
	n := 5
	movies := make([]models.Movie, n)
	for i := 0; i < n; i++ {
		movies[i] = randomMovie()
	}

	testCases := []struct {
		name string
		//setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			/*setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
			},*/
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					ListMovies(gomock.Any()).
					Times(1).
					Return(movies, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				//requireBodyMatchHalls(t, recorder.Body, halls)
			},
		},
		/*{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},*/
		{
			name: "InternalError",
			/*setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
			},*/
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListMovies(gomock.Any()).
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

			//tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestSearchMoviesAPI(t *testing.T) {
	halls := make([]models.Hall, 1)
	halls[0] = randomHall()

	testCases := []struct {
		name string
		//setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			/*setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
			},*/
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetHall(gomock.Any(), gomock.Eq(halls[0].Name)).
					Times(1).
					Return(halls, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchHalls(t, recorder.Body, halls)
			},
		},
		/*{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},*/
		{
			name: "InternalError",
			/*setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
			},*/
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetHall(gomock.Any(), gomock.Eq(halls[0].Name)).
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

			url := "/movies/:id"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			request.URL.RawQuery = q.Encode()

			//tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestInsertMovieAPI(t *testing.T) {
	movie := randomMovie()

	testCases := []struct {
		name string
		body gin.H
		//setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
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
			/*setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
			},*/
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
				store.EXPECT().
					InsertHall(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(&movie, nil)
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
			/*setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
			},*/
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					InsertHall(gomock.Any(), gomock.Any()).
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
			/*setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, user.Role, time.Minute)
			},*/
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					InsertHall(gomock.Any(), gomock.Any()).
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

			//tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateMovieAPI(t *testing.T) {
	movie := randomMovie()

	testCases := []struct {
		name          string
		movieID       string
		body          gin.H
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
			buildStubs: func(store *mockdb.MockStore) {
				arg := models.Movie{
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
				store.EXPECT().
					UpdateMovie(gomock.Any(), gomock.Eq(movie.ID.Hex()), gomock.Eq(arg)).
					Times(1).
					Return(arg, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchMovie(t, recorder.Body, movie)
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
			buildStubs: func(store *mockdb.MockStore) {
				arg := models.Movie{
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
				store.EXPECT().
					UpdateHall(gomock.Any(), gomock.Eq(movie.ID.Hex()), gomock.Eq(arg)).
					Times(1).
					Return(models.Movie{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:    "InvalidRows",
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
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateMovie(gomock.Any(), gomock.Any(), gomock.Any()).
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

			url := "/movies/" + tc.movieID
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteMovieAPI(t *testing.T) {
	movie := randomMovie()
	movieID := movie.ID.Hex()

	testCases := []struct {
		name          string
		movieID       string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			movieID: movieID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteMovie(gomock.Any(), gomock.Eq(movieID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchResponse(t, recorder.Body, apiResponse{Message: "Hall has been deleted"})
			},
		},
		{
			name:    "InternalError",
			movieID: movieID,
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

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomMovie() models.Movie {
	objectID := primitive.NewObjectID()
	return models.Movie{
		ID:        objectID,
		Title:     util.RandomString(50),
		Duration:  int32(util.RandomInt(100, 250)),
		Genre:     util.RandomString(200),
		Directors: util.RandomString(200),
		Actors:    util.RandomString(200),
		Screening: time.Now(),
		Plot:      util.RandomString(200),
		Poster:    util.RandomString(200),
	}
}

func requireBodyMatchMovie(t *testing.T, body *bytes.Buffer, movie models.Movie) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotMovie models.Movie
	err = json.Unmarshal(data, &gotMovie)
	require.NoError(t, err)
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
