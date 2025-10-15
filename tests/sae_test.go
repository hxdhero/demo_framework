package tests

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSAE(t *testing.T) {

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/1/g/1.0/sae/checkpreload/", nil)
	service.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)
	require.Equal(t, "successful", w.Body.String())

}
