package parser_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/marcsanmi/query-benchmarking/internal/parser"
	"github.com/stretchr/testify/assert"
)

// Custom Reader mock.
type readerMock struct {
	buf *bytes.Buffer
}

func newReaderMock(buf bytes.Buffer) readerMock {
	return readerMock{buf: &buf}
}

func (r readerMock) Read(p []byte) (n int, err error) {
	return r.buf.Read(p)
}

func TestParser_ParseQueries(t *testing.T) {
	queriesData := []byte(`demo_cpu_usage_seconds_total{mode="idle"}|1597056698698|1597059548699|15000\n
		demo_cpu_usage_seconds_total{mode="idle"}|1597056698698|1597059548699|15000\n
		invalid_query`)

	readerMock := newReaderMock(createBufferWithContent(queriesData))

	p := parser.New(readerMock)
	q := p.ParseQueries()

	fmt.Printf("queries: %+v \n", q)
	assert.Equal(t, 2, len(q))
}

func createBufferWithContent(input []byte) bytes.Buffer {
	buf := bytes.Buffer{}
	buf.Write(input)
	return buf
}
