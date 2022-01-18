package parser

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/marcsanmi/query-benchmarking/internal"
	log "github.com/sirupsen/logrus"
)

const delimiter string = "|"

type Parser struct {
	file io.Reader
}

func New(file io.Reader) Parser {
	return Parser{file: file}
}

// ParseQueries given a reader, parse its content and returns it in the format of a valid slice of queries.
func (p Parser) ParseQueries() []internal.Query {
	queries := make([]internal.Query, 0)
	s := bufio.NewScanner(p.file)
	for s.Scan() {
		queryStr := strings.TrimSpace(s.Text())
		query, err := p.buildQueryFromString(queryStr)
		if err != nil {
			log.Errorf("skipping row: %v", err.Error())
			continue
		}
		queries = append(queries, query)
	}

	return queries
}

// BuildQueryFromString builds a Query from any given string and adds some validation as well.
func (p Parser) buildQueryFromString(input string) (internal.Query, error) {
	var query internal.Query

	chunks := strings.Split(input, delimiter)
	if len(chunks) != 4 {
		return query, errors.New("wrong input format")
	}

	return internal.Query{
		Params: map[string]string{
			"query": chunks[0],
			"start": parseTimestamp(chunks[1]),
			"end":   parseTimestamp(chunks[2]),
			"step":  chunks[3],
		},
	}, nil
}

// parseTimestamp parse converts the string millisecond timestamps fields into second timestamps.
func parseTimestamp(timestamp string) string {
	if len(timestamp) > 10 {
		timestamp = timestamp[:10] + "." + timestamp[10:]
	}

	return timestamp
}
