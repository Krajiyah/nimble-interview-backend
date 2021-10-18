package users

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Krajiyah/nimble-interview-backend/internal/models"
	"github.com/Krajiyah/nimble-interview-backend/internal/testutils"
	"github.com/Krajiyah/nimble-interview-backend/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestSignupAPIHappyPath(t *testing.T) {
	deps, fileName, err := testutils.NewUnitDeps()
	require.NoError(t, err)
	defer os.Remove(fileName)
	api := NewSignupAPI(deps)

	const (
		username = "testUsername"
		password = "testPassword"
	)
	body := createAuthInput(username, password)

	r, err := http.NewRequest(http.MethodPost, "/users", body)
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	require.NoError(t, api.Handler(c))
	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	responseBody, err := ioutil.ReadAll(w.Result().Body)
	require.NoError(t, err)
	res := utils.Response{}
	require.NoError(t, json.Unmarshal(responseBody, &res))
	token := res.Result.(map[string]interface{})["token"].(string)
	require.NotEmpty(t, token)

	user, err := models.GetUserByUsername(deps.DB(), username)
	require.NoError(t, err)
	require.Equal(t, username, user.Username)
	require.True(t, checkHash(password, user.PasswordHash))

	compareUser, err := utils.ValidateJWT(deps.DB(), token)
	require.NoError(t, err)
	require.Equal(t, user.ID, compareUser.ID)
	require.Equal(t, user.Username, compareUser.Username)
}

func TestSignupAPIInvalidParams(t *testing.T) {
	deps, fileName, err := testutils.NewUnitDeps()
	require.NoError(t, err)
	defer os.Remove(fileName)
	api := NewSignupAPI(deps)

	const (
		username = "testUsername"
		password = ""
	)
	body := createAuthInput(username, password)

	r, err := http.NewRequest(http.MethodPost, "/users", body)
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	require.NoError(t, api.Handler(c))
	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	user, err := models.GetUserByUsername(deps.DB(), username)
	require.Nil(t, user)
	require.Error(t, err)
}

func TestSignupAlreadyExists(t *testing.T) {
	deps, fileName, err := testutils.NewUnitDeps()
	require.NoError(t, err)
	defer os.Remove(fileName)
	api := NewSignupAPI(deps)

	const (
		username = "testUsername"
		password = "testPassword"
	)
	body := createAuthInput(username, password)

	r, err := http.NewRequest(http.MethodPost, "/users", body)
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	require.NoError(t, api.Handler(c))
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	user, err := models.GetUserByUsername(deps.DB(), username)
	require.NoError(t, err)
	require.Equal(t, username, user.Username)
	require.True(t, checkHash(password, user.PasswordHash))

	w = httptest.NewRecorder()
	c = echo.New().NewContext(r, w)
	require.NoError(t, api.Handler(c))
	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestLoginAPIHappyPath(t *testing.T) {
	deps, fileName, err := testutils.NewUnitDeps()
	require.NoError(t, err)
	defer os.Remove(fileName)
	api := NewSignupAPI(deps)

	const (
		username = "testUsername"
		password = "testPassword"
	)
	body := createAuthInput(username, password)

	r, err := http.NewRequest(http.MethodPost, "/users", body)
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	require.NoError(t, api.Handler(c))
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	user, err := models.GetUserByUsername(deps.DB(), username)
	require.NoError(t, err)
	require.Equal(t, username, user.Username)
	require.True(t, checkHash(password, user.PasswordHash))

	body = createAuthInput(username, password)
	r, err = http.NewRequest(http.MethodPost, "/users/login", body)
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	require.NoError(t, err)
	w = httptest.NewRecorder()
	c = echo.New().NewContext(r, w)
	require.NoError(t, NewLoginAPI(deps).Handler(c))
	require.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestLoginAPIInvalidPassword(t *testing.T) {
	deps, fileName, err := testutils.NewUnitDeps()
	require.NoError(t, err)
	defer os.Remove(fileName)
	api := NewSignupAPI(deps)

	const (
		username = "testUsername"
		password = "testPassword"
	)
	body := createAuthInput(username, password)

	r, err := http.NewRequest(http.MethodPost, "/users", body)
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	c := echo.New().NewContext(r, w)
	require.NoError(t, api.Handler(c))
	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	user, err := models.GetUserByUsername(deps.DB(), username)
	require.NoError(t, err)
	require.Equal(t, username, user.Username)
	require.True(t, checkHash(password, user.PasswordHash))

	body = createAuthInput(username, "badPassword")
	r, err = http.NewRequest(http.MethodPost, "/users/login", body)
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	require.NoError(t, err)
	w = httptest.NewRecorder()
	c = echo.New().NewContext(r, w)
	require.NoError(t, NewLoginAPI(deps).Handler(c))
	require.Equal(t, http.StatusForbidden, w.Result().StatusCode)
}

func createAuthInput(username, password string) io.Reader {
	return strings.NewReader(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password))
}
