package messages

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
	"github.com/Krajiyah/nimble-interview-backend/internal/routes/middlewares"
	"github.com/Krajiyah/nimble-interview-backend/internal/testutils"
	"github.com/Krajiyah/nimble-interview-backend/internal/utils"
	mapset "github.com/deckarep/golang-set"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestMessageAPIs(t *testing.T) {
	deps, fileName, err := testutils.NewUnitDeps()
	require.NoError(t, err)
	defer os.Remove(fileName)
	sendMessageAPI := NewSendMessageAPI(deps)
	getMessagesAPI := NewGetMessagesAPI(deps)

	user, err := models.NewUser(deps.DB(), &models.User{
		Username:     "someUserName",
		PasswordHash: "someHashOfPassword",
	})
	require.NoError(t, err)

	const (
		numMessages = 100
		pageSize    = 10
	)

	i := 0
	for i < numMessages {
		body := createMessageInput(fmt.Sprintf("some message data %d", i))
		r, err := http.NewRequest(http.MethodPost, "/messages", body)
		r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		require.NoError(t, err)
		w := httptest.NewRecorder()
		c := echo.New().NewContext(r, w)
		c.Set(middlewares.UserContextKey, user)
		require.NoError(t, sendMessageAPI.Handler(c))
		require.Equal(t, http.StatusOK, w.Result().StatusCode)
		i++
	}

	page := 1
	messages := mapset.NewSet()
	for page <= numMessages/pageSize {
		r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/messages?page=%d&pageSize=%d", page, pageSize), nil)
		require.NoError(t, err)
		w := httptest.NewRecorder()
		c := echo.New().NewContext(r, w)
		c.Set(middlewares.UserContextKey, user)
		require.NoError(t, getMessagesAPI.Handler(c))
		require.Equal(t, http.StatusOK, w.Result().StatusCode)
		body, err := ioutil.ReadAll(w.Result().Body)
		require.NoError(t, err)
		res := utils.Response{}
		require.NoError(t, json.Unmarshal(body, &res), string(body))
		pageMessages := res.Result.([]interface{})
		for _, messageInt := range pageMessages {
			message := messageInt.(map[string]interface{})
			messages.Add(message["ID"])
		}
		page++
	}

	require.Equal(t, numMessages, messages.Cardinality())
}

func createMessageInput(data string) io.Reader {
	return strings.NewReader(fmt.Sprintf(`{"data": "%s"}`, data))
}
