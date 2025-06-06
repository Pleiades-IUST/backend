package auth

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Pleiades-IUST/backend/utils/config"
	"github.com/Pleiades-IUST/backend/utils/dbutil"
	"github.com/Pleiades-IUST/backend/webservice/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func protected(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func signup(ctx *gin.Context) {
	tx := dbutil.GormDB(ctx)

	var req SignUpRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	err := validateSignupRequest(tx, req)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	bcryptPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	user := user.User{
		Email:    req.Email,
		Username: req.Username,
		Password: string(bcryptPass),
	}

	err = tx.Create(&user).Error
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// generate Token
	tokenString, err := createLoginToken(strconv.FormatInt(user.ID, 10), time.Hour*12, []byte(config.GetSecretKey()))
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	response := SignUpResponse{
		Token: tokenString,
	}

	ctx.JSON(http.StatusOK, response)
}

func login(ctx *gin.Context) {
	db := dbutil.GormDB(ctx.Request.Context())

	var req LoginRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	user, err := user.FetchUserByUsername(db, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.String(http.StatusBadRequest, "username or password is incorrect")
			return
		} else {
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		ctx.Status(http.StatusBadRequest)
		return
	} else if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// generate Token
	tokenString, err := createLoginToken(strconv.FormatInt(user.ID, 10), time.Hour*12, []byte(config.GetSecretKey()))
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Token: tokenString,
	}

	ctx.JSON(http.StatusOK, response)
}

func createLoginToken(userID string, duration time.Duration, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		Issuer:    "IdeYar",
		Subject:   userID,
		Audience:  jwt.ClaimStrings{"Login"},
	})

	return token.SignedString(key)
}

func validateSignupRequest(tx *gorm.DB, req SignUpRequest) error {
	if !ValidateEmail(req.Email) {
		return errors.New("invalid email")
	}

	if !ValidateUsername(req.Username) {
		return errors.New("invalid username")
	}

	if !ValidatePassword(req.Password) {
		return errors.New("invalid password")
	}

	if isEmailDuplicate(tx, req) {
		return errors.New("email already exists")
	}

	if isUsernameDuplicate(tx, req) {
		return errors.New("username already exists")
	}

	return nil
}

func isEmailDuplicate(tx *gorm.DB, req SignUpRequest) bool {
	var duplicateEmail bool
	_ = tx.Raw(`
		SELECT count(*) > 0
		FROM user_t
		WHERE email = ?
	`, req.Email).Scan(&duplicateEmail).Error

	return duplicateEmail
}

func isUsernameDuplicate(tx *gorm.DB, req SignUpRequest) bool {
	var duplicateUsername bool
	_ = tx.Raw(`
		SELECT count(*) > 0
		FROM user_t
		WHERE username = ?
	`, req.Username).Scan(&duplicateUsername).Error

	return duplicateUsername
}
