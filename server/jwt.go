package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AsishMandoi/iitk-coin/database"
	"github.com/AsishMandoi/iitk-coin/global"
	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/nu7hatch/gouuid"
)

// Generate a JWT for a given user (rollno, batch, role)
func GenJWT(usr struct {
	Rollno int
	Email  string
	Batch  string
	Role   string
}) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	// JWT id stores a uuid. This ensures that each JWT is unique, which can help prevent the JWT from being replayed.
	// However, I have used jwtId in the signature of the JWT, this will make sure that once a new JWT is generated,
	// the older one will not work.
	// The jwtId is stored in the redis db with the same expiration period as the JWT itslef, and can be retrieved
	// easily for the verification of the JWT.
	uuid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	jwtId := uuid.String()

	claims := token.Claims.(jwt.MapClaims)

	// The rollno, batch, role ("Admin"/"CoreTeam"/"Coordinator"/nil) and the expiration time are used as claims
	claims["rollno"] = usr.Rollno
	claims["batch"] = usr.Batch
	claims["role"] = usr.Role

	// Change "usr.Email" to "global.MyGmailId" for testing purposes
	claims["email"] = usr.Email
	claims["exp"] = time.Now().Add(time.Duration(global.TknExpTime) * time.Minute).Unix()

	// A token (tokenString) is generated using the header, the payload (using the claims) and the encoded secret key (gobal.SignatureKey)
	tokenString, err := token.SignedString([]byte(global.SignatureKey + jwtId))
	if err != nil {
		return "", err
	}

	err = database.SetJWTid(usr.Rollno, jwtId)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// This authorizes a request (with a token) and returns the claims if successful else returns errors and a status code
func ValidateJWT(r *http.Request) (int, jwt.MapClaims, error) {

	tokenString := r.Header.Get("Authorization")

	// If the header `Authorization` is empty...
	if tokenString = strings.TrimPrefix(tokenString, "Bearer "); tokenString == "" {
		return 400, nil, fmt.Errorf("Authorization token not found")
	}

	// Get the (unverified) token from the header.
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != "HS256" {
			return []byte(""), fmt.Errorf("Invalid signing method")
		}
		rollno := int(t.Claims.(jwt.MapClaims)["rollno"].(float64))

		// Getting the JWT id from an unverified token
		jwtId, err := database.GetJWTid(rollno)
		if err != nil {
			return []byte(""), err
		}
		return []byte(global.SignatureKey + jwtId), nil
	})

	if err != nil {
		return 400, nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 401, nil, fmt.Errorf("Invalid authorization token")
	}
	return 200, claims, nil
}
