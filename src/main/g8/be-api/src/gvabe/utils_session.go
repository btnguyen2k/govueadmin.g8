package gvabe

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"main/src/utils"
)

var (
	errorInvalidClient = errors.New("invalid client id")
	errorInvalidJwt    = errors.New("cannot decode token")
	errorExpiredJwt    = errors.New("token has expired")
)

const (
	loginChannelForm  = "form"
	loginChannelExter = "exter"
)

// Session captures a user-login-session. Session object is to be serialized and embedded into a SessionClaims.
// available since template-v0.2.0
type Session struct {
	ClientRef   string    `json:"cref"` // reference to client data
	Channel     string    `json:"chan"` // login source (form-based, Exter, etc)
	UserId      string    `json:"uid"`  // id of logged-in user
	DisplayName string    `json:"name"` // display name of logged-in user
	CreatedAt   time.Time `json:"cat"`  // timestamp when the session is created
	ExpiredAt   time.Time `json:"eat"`  // timestamp when the session expires
	Data        []byte    `json:"data"` // session's arbitrary data
}

// SessionClaims is an extended structure of JWT's standard claims
// available since template-v0.2.0
type SessionClaims struct {
	UserId          string `json:"uid,omitempty"`  // id of logged-in user
	UserDisplayName string `json:"name,omitempty"` // display name of logged-in user
	Data            []byte `json:"data,omitempty"` // session's arbitrary data
	jwt.StandardClaims
}

func (s *SessionClaims) isExpired() bool {
	return s.ExpiresAt > 0 && s.ExpiresAt < time.Now().Unix()
}

func (s *SessionClaims) isGoingExpired(numSec int64) bool {
	return s.ExpiresAt > 0 && s.ExpiresAt-numSec < time.Now().Unix()
}

/*----------------------------------------------------------------------*/


// available since template-v0.2.0
func parseJwt(jwtStr string, pubKey *rsa.PublicKey) (map[string]interface{}, error) {
	token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("enexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid claim")
	}
}

// available since template-v0.2.0
func genJws(claim *SessionClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	return token.SignedString(rsaPrivKey)
}

// available since template-v0.2.0
func parseLoginToken(jwtStr string) (*SessionClaims, error) {
	claims, err := parseJwt(jwtStr, rsaPubKey)
	if err != nil {
		return nil, err
	}
	var result SessionClaims
	js, _ := json.Marshal(claims)
	return &result, json.Unmarshal(js, &result)
}

// genLoginClaims generates a login token as SessionClaims:
//
// available since template-v0.2.0
func genLoginClaims(id string, sess *Session) (*SessionClaims, error) {
	if id == "" {
		id = utils.UniqueId()
	}
	u, err := userDaov2.Get(sess.UserId)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New(fmt.Sprintf("user [%s] not found", sess.UserId))
	}
	sessData, err := json.Marshal(sess)
	if err != nil {
		return nil, err
	}
	sessData, err = zipAndEncrypt(sessData)
	return &SessionClaims{
		UserId:          sess.UserId,
		UserDisplayName: sess.DisplayName,
		Data:            sessData,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: sess.ExpiredAt.Unix(),
			Id:        id,
			IssuedAt:  sess.CreatedAt.Unix(),
			Subject:   sess.Channel,
		},
	}, err
}
