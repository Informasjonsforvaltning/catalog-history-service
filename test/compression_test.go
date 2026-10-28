package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Informasjonsforvaltning/catalog-history-service/config"
	"github.com/stretchr/testify/assert"
)

func getUpdatesWithAcceptEncoding(t *testing.T, acceptEncoding string) *httptest.ResponseRecorder {
	router := config.SetupRouter()

	jwt := CreateMockJwt(time.Now().Add(time.Hour).Unix(), &TestValues.SysAdminAuth, &TestValues.Audience)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/111222333/123456789/updates", nil)
	req.Header.Set("Authorization", *jwt)
	if acceptEncoding != "" {
		req.Header.Set("Accept-Encoding", acceptEncoding)
	}
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	return w
}

func TestGzipsResponseWhenClientAcceptsIt(t *testing.T) {
	w := getUpdatesWithAcceptEncoding(t, "gzip")

	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
	assert.Contains(t, w.Header().Get("Vary"), "Accept-Encoding")

	body := w.Body.Bytes()
	assert.True(t, len(body) > 2 && body[0] == 0x1f && body[1] == 0x8b, "body is not gzip")
}

func TestLeavesResponseUncompressedWhenClientDoesNotAsk(t *testing.T) {
	w := getUpdatesWithAcceptEncoding(t, "")

	assert.Empty(t, w.Header().Get("Content-Encoding"))
	assert.True(t, w.Body.Len() > 0, "response body must not be empty")
}
