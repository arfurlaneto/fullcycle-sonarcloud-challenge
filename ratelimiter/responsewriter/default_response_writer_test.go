package responsewriter

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NewRateLimiterDefaultResponseWriterTestSuite struct {
	suite.Suite
}

func TestNewRateLimiterDefaultResponseWriterTestSuite(t *testing.T) {
	suite.Run(t, new(NewRateLimiterDefaultResponseWriterTestSuite))
}

func (s *NewRateLimiterDefaultResponseWriterTestSuite) TestWriteResponse() {
	recorder := httptest.NewRecorder()
	writer := http.ResponseWriter(recorder)

	responseWriter := NewRateLimiterDefaultResponseWriter()
	responseWriter.WriteResponse(&writer)

	response := recorder.Result()
	responseStatus := response.StatusCode
	responseBody, err := ioutil.ReadAll(response.Body)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 429, responseStatus)
	assert.Equal(s.T(), "you have reached the maximum number of requests or actions allowed within a certain time frame", string(responseBody))

}

func (s *NewRateLimiterDefaultResponseWriterTestSuite) TestWriteError() {
	recorder := httptest.NewRecorder()
	writer := http.ResponseWriter(recorder)

	responseWriter := NewRateLimiterDefaultResponseWriter()
	responseWriter.WriteError(&writer, errors.New("error"))

	response := recorder.Result()
	responseStatus := response.StatusCode
	responseBody, err := ioutil.ReadAll(response.Body)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 500, responseStatus)
	assert.Equal(s.T(), "internal server error", string(responseBody))
}
