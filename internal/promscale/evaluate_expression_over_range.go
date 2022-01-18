package promscale

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/marcsanmi/query-benchmarking/internal"
)

const getRangeExpressionEndpoint string = "/api/v1/query_range"

type EvaluateExpressionQueryOverRangeResponse struct {
	Status string `json:"status"`
}

// EvaluateExpressionQueryOverRange evaluates an expression query over a range of time
func (c Client) EvaluateExpressionQueryOverRange(ctx context.Context, query internal.Query) (elapsed time.Duration, err error) {
	defer func(start time.Time) {
		elapsed = time.Since(start)
	}(time.Now())

	path, err := c.baseURL.Parse(getRangeExpressionEndpoint)
	if err != nil {
		return elapsed, fmt.Errorf("create path for request: %v", err)
	}

	q := path.Query()
	for k, v := range query.Params {
		q.Set(k, v)
	}
	path.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path.String(), nil)
	if err != nil {
		return elapsed, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return elapsed, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return elapsed, fmt.Errorf("wrong http status code received: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return elapsed, err
	}

	var expressionRes EvaluateExpressionQueryOverRangeResponse
	err = json.Unmarshal(body, &expressionRes)
	if err != nil {
		return elapsed, err
	}

	// Check response body status is success.
	if expressionRes.Status != "success" {
		return elapsed, fmt.Errorf("wrong response body status received: %v", expressionRes.Status)
	}

	return elapsed, nil
}
