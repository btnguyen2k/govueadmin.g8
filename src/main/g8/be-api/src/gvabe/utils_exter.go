package gvabe

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/consu/semita"

	"main/src/goapi"
	userv2 "main/src/gvabe/bov2/user"
	"main/src/utils"
)

var (
	exterClient    *ExterClient
	exterRsaPubKey *rsa.PublicKey
)

// available since template-v0.2.0
func NewExterClient(appId, baseUrl string) *ExterClient {
	return &ExterClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		appId:      appId,
		baseUrl:    baseUrl,
	}
}

// available since template-v0.2.0
type ExterClient struct {
	httpClient *http.Client
	appId      string
	baseUrl    string
}

func (ec *ExterClient) parseExterResponse(resp *http.Response) (*ExterResponse, error) {
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("http response status %d", resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	eresp := ExterResponse{
		raw: make(map[string]interface{}),
	}
	if err := json.Unmarshal(body, &(eresp.raw)); err != nil {
		return nil, err
	}
	eresp.s = semita.NewSemita(eresp.raw)
	eresp.Status = eresp.GetInt("status")
	eresp.Message = eresp.GetString("message")
	return &eresp, nil
}

func (ec *ExterClient) doRequest(method, apiUri string, body []byte) (*ExterResponse, error) {
	var reader io.Reader = nil
	if body != nil {
		reader = bytes.NewBuffer(body)
	}
	req, _ := http.NewRequest(method, ec.baseUrl+apiUri, reader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Id", ec.appId)
	if resp, err := ec.httpClient.Do(req); err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
		return ec.parseExterResponse(resp)
	}
}

func (ec *ExterClient) Info() (*ExterResponse, error) {
	return ec.doRequest("GET", "/info", nil)
}

func (ec *ExterClient) VerifyLoginToken(token string) (*ExterResponse, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{"token": token, "app": ec.appId})
	return ec.doRequest("POST", "/api/verifyLoginToken", reqBody)
}

/*----------------------------------------------------------------------*/

// available since template-v0.2.0
type ExterResponse struct {
	Status  int
	Message string
	raw     map[string]interface{}
	s       *semita.Semita
}

func (er *ExterResponse) GetDataAsType(path string, typ reflect.Type) (interface{}, error) {
	return er.s.GetValueOfType(path, typ)
}

func (er *ExterResponse) GetString(path string) string {
	if v, err := er.GetDataAsType(path, reddo.TypeString); err == nil && v != nil {
		return v.(string)
	}
	return ""
}

func (er *ExterResponse) GetInt(path string) int {
	if v, err := er.GetDataAsType(path, reddo.TypeInt); err == nil && v != nil {
		return int(v.(int64))
	}
	return 0
}

/*----------------------------------------------------------------------*/

// available since template-v0.2.0
func goFetchExterInfo(sleepSeconds int) {
	if sleepSeconds < 60 {
		sleepSeconds = 60
	}
	for ; ; {
		resp, err := exterClient.Info()
		if err != nil {
			log.Printf("[ERROR] goFetchExterInfo - Error calling Exter api: 0/%s", err)
		} else if resp.Status == 200 {
			pubKeyPem := resp.GetString("data.rsa_public_key")
			pubKey, err := parseRsaPublicKeyFromPem(pubKeyPem)
			if err != nil {
				log.Printf("[ERROR] goFetchExterInfo - Cannot extract Exter RSA public key: %e / %v", err, resp.raw)
			} else {
				exterRsaPubKey = pubKey
				log.Printf("[DEBUG_MODE] Exter public key: {Size: %d / Exponent: %d / Modulus: %x}",
					exterRsaPubKey.Size()*8, exterRsaPubKey.E, exterRsaPubKey.N)
			}
		} else {
			log.Printf("[ERROR] goFetchExterInfo - Error calling Exter api: %d / %s", resp.Status, resp.Message)
		}
		time.Sleep(time.Duration(sleepSeconds) * time.Second)
	}
}

// available since template-v0.2.0
type ExterToken struct {
	Id        string `json:"jti,omitempty"`  // token's unique id
	AppId     string `json:"aud,omitempty"`  // client app id
	Type      string `json:"type,omitempty"` // token type
	Data      []byte `json:"data,omitempty"` // internal use by Exter
	ExpiresAt int64  `json:"exp,omitempty"`  // expiry UNIX-timestamp
	IssuedAt  int64  `json:"iat,omitempty"`  // issue UNIX-timestamp
	UserName  string `json:"name,omitempty"` // user's display name
	UserId    string `json:"uid,omitempty"`  // user's id
	Channel   string `json:"sub,omitempty"`  // login channel / identity source
}

// available since template-v0.2.0
func parseExterJwt(jwtStr string) (*ExterToken, error) {
	jwtData, err := parseJwt(jwtStr, exterRsaPubKey)
	if err != nil || jwtData == nil {
		return nil, err
	}
	js, _ := json.Marshal(jwtData)
	result := ExterToken{}
	return &result, json.Unmarshal(js, &result)
}

// available since template-v0.2.0
func createUserFromExterToken(exterToken *ExterToken) (*userv2.User, error) {
	if exterToken == nil {
		return nil, nil
	}
	if exterToken.UserId == "" {
		return nil, errors.New("no user-id found in Exter token")
	}
	user, err := userDaov2.Get(exterToken.UserId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error while getting user [%s]: %e", exterToken.UserId, err))
	}
	if user == nil {
		log.Printf("[INFO] Creating user [%s] from Exter token...", exterToken.UserId)
		user = userv2.NewUser(goapi.AppVersionNumber, exterToken.UserId, utils.UniqueId())
		displayName := exterToken.UserName
		if displayName == "" {
			displayName = user.GetMaskId()
		}
		user.SetDisplayName(displayName).SetAdmin(false)
		_, err = userDaov2.Create(user)
	}
	return user, err
}
