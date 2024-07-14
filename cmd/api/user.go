package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tijanadmi/movieginmongoapi/models"
	"github.com/tijanadmi/movieginmongoapi/util"
	"go.mongodb.org/mongo-driver/mongo"
)

func newUserResponse(user *models.User) userResponse {
	return userResponse{
		Username: user.Username,
	}
}

// Paths Information

// @Summary Provides a JSON Web Token
// @Description Authenticates a user and provides a Paseto/JWT to Authorize API calls
// @ID loginUser
// @Consume json
// @Produce json
// @Param loginUserRequest body loginUserRequest true "User login request"
// @Success 200 {object} loginUserResponse
// @Failure 400 {object} apiErrorResponse
// @Failure 401 {object} apiErrorResponse
// @Router /users/login [post]
func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apiErrorResponse{Error: err.Error()})
		return
	}

	if !util.IsValidUsername(req.Username) {
		ctx.JSON(http.StatusBadRequest, apiErrorResponse{Error: "invalid username"})
		return
	}


	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, apiErrorResponse{Error: "user not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		}
		return
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, apiErrorResponse{Error: err.Error()})
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		/*user.Role,*/
		"user",
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		//user.Role,
		"user",
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	rsp := loginUserResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getUserByUsername(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, apiErrorResponse{Error: err.Error()})
		return
	}

	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, apiErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (server *Server) InsertUser(ctx *gin.Context) {
	var user *models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"invalid input": err.Error(),
		})
		return
	}

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user.Password = hashedPassword
	/*user1, _ := server.store.GetUserByUsername(ctx, user.Username)

	if user.Username == user1.Username {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "User already exists",
		})
		return
	}*/

	if user, err = server.store.InsertUser(ctx, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}
