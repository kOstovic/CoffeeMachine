package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kOstovic/CoffeeMachine/cmd/CoffeeMachineRestApiV3/config"
	log "github.com/sirupsen/logrus"
)

var jwtKey = []byte("my_key")
var tokens []string

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func RegisterRoutesAuth(router *gin.RouterGroup) {

	router.POST("/", gin.BasicAuth(gin.Accounts{
		config.Configuration.Auth.USERNAME: config.Configuration.Auth.PASSWORD,
	}), tokenEndpoint)
}

// tokenEndpoint godoc
// @Summary Login to administrator CoffeeMachine
// @Description Login to administrator CoffeeMachine
// @Produce json
// @Success 200 {object} bool
// @Failure 400,401,404
// @Failure 500
// @securityDefinitions.http BasicAuth
// @in header
// @name Authorization
// @Router /login [post]
func tokenEndpoint(c *gin.Context) {
	user, userExist := c.Get(gin.AuthUserKey)
	if !userExist {
		log.Errorf("Could not generate JWT")
	}
	token, err := generateJWT(user.(string))
	if err != nil {
		log.Errorf("Could not generate JWT")
	}
	tokens = append(tokens, token)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func parseToken(c *gin.Context) (int, error) {
	bearerToken := c.Request.Header.Get("Authorization")
	if bearerToken == "" {
		log.Warnf("No Bearer token provided")
		return http.StatusUnauthorized, fmt.Errorf("No Bearer token provided")
	}
	parsedBearer := strings.Split(bearerToken, " ")
	if len(parsedBearer) <= 1 {
		log.Warnf("Malformed Bearer token provided")
		return http.StatusUnauthorized, fmt.Errorf("Malformed Bearer token provided")
	}
	reqToken := parsedBearer[1]
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Warnf("Invalid signature in jwt bearer token: %v", bearerToken)
			return http.StatusUnauthorized, fmt.Errorf("message: unauthorized")
		}
		log.Warnf("Error in token: %v: for Bearer token: %v", tkn, bearerToken)
		return http.StatusBadRequest, fmt.Errorf("message: bad request")
	}
	if !tkn.Valid {
		log.Warnf("Invalid token: %v: for Bearer token: %v", tkn, bearerToken)
		return http.StatusUnauthorized, fmt.Errorf("message: unauthorized")
	}
	return 0, nil
}

func generateJWT(user string) (string, error) {
	log.Debugf("Generating JWT")
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		Username: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}
