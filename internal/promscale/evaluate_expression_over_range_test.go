package promscale_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marcsanmi/query-benchmarking/internal"
	"github.com/marcsanmi/query-benchmarking/internal/promscale"
	"github.com/stretchr/testify/assert"
)

const evaluateExpressionResponseBody = `{"status": "success"}`

func TestClient_EvaluateExpressionOverRange(t *testing.T) {
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(evaluateExpressionResponseBody))
		assert.Equal(t, r.Method, "GET")
	}))
	defer httpServer.Close()

	client, err := promscale.NewClient(&http.Client{}, promscale.Config{URL: httpServer.URL})
	assert.NoError(t, err)

	query := internal.Query{Params: map[string]string{"start": "1123456789"}}
	elapsed, err := client.EvaluateExpressionQueryOverRange(context.Background(), query)
	assert.NoError(t, err)
	assert.NotEmpty(t, elapsed)
}
