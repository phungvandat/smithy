package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/go-kit/kit/log"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	auth "github.com/dwarvesf/smithy/backend/auth"
	"github.com/dwarvesf/smithy/backend/endpoints"
	"github.com/dwarvesf/smithy/backend/service"
	utilTest "github.com/dwarvesf/smithy/common/utils/database/pg/test/set1"
)

const (
	secretKey string = "lalala"
	Admin     string = "admin"
	User      string = "user"
)

func TestNewHTTPHandler(t *testing.T) {
	//make up-dashboard
	tsDashboard := httptest.NewServer(initDashboardServer(t))
	defer tsDashboard.Close()

	tests := []struct {
		name       string
		header     http.Header
		wantErr    string
		wantStatus int
	}{
		{
			name:       "Success",
			header:     newAuthHeader(auth.New(secretKey, "aaa", Admin).Encode()),
			wantStatus: http.StatusOK,
		},
		{
			name:       "Header is nil",
			header:     nil,
			wantErr:    "jwtauth: no token found",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Header is empty",
			header:     http.Header{},
			wantErr:    "jwtauth: no token found",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "wrong secret key",
			header:     newAuthHeader(auth.New("wrong", "aaa", Admin).Encode()),
			wantErr:    "signature is invalid",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "wrong token string",
			header:     newAuthHeader("blabla"),
			wantErr:    "token contains an invalid number of segments",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "wrong algorithm",
			header: newAuthHeader(newJwt512Token([]byte(secretKey), jwtauth.Claims{
				"username": "aaa",
				"role":     Admin,
			})),
			wantErr:    "jwtauth: token is unauthorized",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "wrong secret key and algorithm",
			header: newAuthHeader(newJwt512Token([]byte("wrong"), jwtauth.Claims{
				"username": "aaa",
				"role":     Admin,
			})),
			wantErr:    "signature is invalid",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Client can't agent-sync",
			header:     newAuthHeader(auth.New(secretKey, "bbb", User).Encode()),
			wantErr:    "Unauthorized",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if status, resp := testRequest(t, tsDashboard, "GET", "/models", tt.header, nil); status != tt.wantStatus || resp != tt.wantErr {
				t.Errorf("NewHTTPHandler() = (%v, %d), want (%v, %d)", resp, status, tt.wantErr, tt.wantStatus)
			}
		})
	}

	// test login api
	loginTests := []struct {
		name       string
		jsonString []byte
		wantStatus int
		wantErr    bool
	}{
		{
			name:       "Login success",
			jsonString: []byte(`{"username":"aaa", "password": "abc"}`),
			wantStatus: 200,
			wantErr:    false,
		},
		{
			name:       "Login wrong username",
			jsonString: []byte(`{"username":"adfs", "password": "abc"}`),
			wantStatus: 200,
			wantErr:    false,
		},
		{
			name:       "Login wrong password",
			jsonString: []byte(`{"username":"aaa", "password": "conmeocon"}`),
			wantStatus: 401,
			wantErr:    true,
		},
	}

	for _, tt := range loginTests {
		t.Run(tt.name, func(t *testing.T) {
			if status := loginTestRequest(t, tsDashboard, "/auth/login", tt.jsonString); status != tt.wantStatus && tt.wantErr {
				t.Errorf("Login() = %v, want %v", status, tt.wantStatus)
			}
		})
	}
}

//
// Test helper functions
//

func testRequest(t *testing.T, ts *httptest.Server, method, path string, header http.Header, body io.Reader) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return 0, ""
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v[0])
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return 0, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return 0, ""
	}
	defer resp.Body.Close()

	data := &auth.ErrAuthentication{}

	if err = json.Unmarshal(respBody, data); err != nil {
		return resp.StatusCode, string(respBody)
	}

	return resp.StatusCode, data.Error
}

func newAuthHeader(tokenStr string) http.Header {
	h := http.Header{}
	h.Set("Authorization", "BEARER "+tokenStr)
	return h
}

func initDashboardServer(t *testing.T) http.Handler {
	cfg, _ := utilTest.CreateConfig(t)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	s, err := service.NewService(cfg)
	if err != nil {
		t.Fatal(err)
	}

	return NewHTTPHandler(
		endpoints.MakeServerEndpoints(s),
		logger,
		os.Getenv("ENV") == "local",
		cfg.Authentication.SerectKey,
	)
}

func newJwt512Token(secret []byte, claims ...jwtauth.Claims) string {
	// use-case: when token is signed with a different alg than expected
	token := jwt.New(jwt.GetSigningMethod("HS512"))
	if len(claims) > 0 {
		token.Claims = claims[0]
	}
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		fmt.Println("error at newJwt512Token")
	}
	return tokenStr
}

func loginTestRequest(t *testing.T, ts *httptest.Server, path string, body []byte) int {
	req, err := http.NewRequest("POST", ts.URL+path, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return 0
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, errRespone := client.Do(req)
	if errRespone != nil {
		t.Fatal(err)
		return 0
	}

	return response.StatusCode
}
