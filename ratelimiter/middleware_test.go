package ratelimiter

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/arfurlaneto/fullcycle-sonarcloud-challenge/ratelimiter/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type MiddlewareTestSuite struct {
	suite.Suite
	controller         *gomock.Controller
	context            context.Context
	responseWriterMock *mocks.MockRateLimiterResponseWriter
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}

func (s *MiddlewareTestSuite) SetupTest() {
	s.controller = gomock.NewController(s.T())
	s.context = context.Background()
	s.responseWriterMock = mocks.NewMockRateLimiterResponseWriter(s.controller)
}

func (s *MiddlewareTestSuite) TestMiddleware_NewRateLimiter() {
	middleware := NewRateLimiter()
	assert.NotNil(s.T(), middleware)
}

func (s *MiddlewareTestSuite) TestMiddleware_NewRateLimiterWithConfig() {
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := NewRateLimiter()
	middlewareFunc := middleware(emptyHandler)
	assert.NotNil(s.T(), middleware)
	assert.NotNil(s.T(), middlewareFunc)
}

func (s *MiddlewareTestSuite) TestMiddleware_NewRateLimiterWithConfig2() {
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := NewRateLimiterWithConfig(&RateLimiterConfig{})
	middlewareFunc := middleware(emptyHandler)
	assert.NotNil(s.T(), middleware)
	assert.NotNil(s.T(), middlewareFunc)
}

func (s *MiddlewareTestSuite) TestMiddleware_IPAllowed() {
	config := &RateLimiterConfig{
		IP: &RateLimiterRateConfig{
			MaxRequestsPerSecond:  10,
			BlockTimeMilliseconds: 100,
		},
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("DONE"))
	})

	rateLimiterCheckFunction := func(ctx context.Context, keyType string, key string, config *RateLimiterConfig, rateConfig *RateLimiterRateConfig) (*time.Time, error) {
		return nil, nil
	}

	request := httptest.NewRequest("GET", "http://testing", nil)
	recorder := httptest.NewRecorder()

	rateLimiter(config, nextHandler, rateLimiterCheckFunction).ServeHTTP(recorder, request)

	response := recorder.Result()
	responseStatus := response.StatusCode
	responseBody, err := ioutil.ReadAll(response.Body)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 200, responseStatus)
	assert.Equal(s.T(), "DONE", string(responseBody))
}

func (s *MiddlewareTestSuite) TestMiddleware_IPNotAllowed() {
	config := &RateLimiterConfig{
		IP: &RateLimiterRateConfig{
			MaxRequestsPerSecond:  10,
			BlockTimeMilliseconds: 100,
		},
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("DONE"))
	})

	rateLimiterCheckFunction := func(ctx context.Context, keyType string, key string, config *RateLimiterConfig, rateConfig *RateLimiterRateConfig) (*time.Time, error) {
		block := time.Now().Add(time.Millisecond * 100)
		return &block, nil
	}

	request := httptest.NewRequest("GET", "http://testing", nil)
	recorder := httptest.NewRecorder()

	s.responseWriterMock.EXPECT().WriteResponse(gomock.Any()).Do(func(w *http.ResponseWriter) {
		(*w).WriteHeader(429)
		(*w).Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
	})
	config.ResponseWriter = s.responseWriterMock

	rateLimiter(config, nextHandler, rateLimiterCheckFunction).ServeHTTP(recorder, request)

	response := recorder.Result()
	responseStatus := response.StatusCode
	responseBody, err := ioutil.ReadAll(response.Body)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 429, responseStatus)
	assert.Equal(s.T(), "you have reached the maximum number of requests or actions allowed within a certain time frame", string(responseBody))
}

func (s *MiddlewareTestSuite) TestMiddleware_IPError() {
	config := &RateLimiterConfig{
		IP: &RateLimiterRateConfig{
			MaxRequestsPerSecond:  10,
			BlockTimeMilliseconds: 100,
		},
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("DONE"))
	})

	rateLimiterCheckFunction := func(ctx context.Context, keyType string, key string, config *RateLimiterConfig, rateConfig *RateLimiterRateConfig) (*time.Time, error) {
		return nil, errors.New("error")
	}

	request := httptest.NewRequest("GET", "http://testing", nil)
	recorder := httptest.NewRecorder()

	s.responseWriterMock.EXPECT().WriteError(gomock.Any(), gomock.Any()).Do(func(w *http.ResponseWriter, err error) {
		(*w).WriteHeader(500)
		(*w).Write([]byte("internal server error"))
	})
	config.ResponseWriter = s.responseWriterMock

	rateLimiter(config, nextHandler, rateLimiterCheckFunction).ServeHTTP(recorder, request)

	response := recorder.Result()
	responseStatus := response.StatusCode
	responseBody, err := ioutil.ReadAll(response.Body)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 500, responseStatus)
	assert.Equal(s.T(), "internal server error", string(responseBody))
}

func (s *MiddlewareTestSuite) TestMiddleware_TokenAllowed() {
	config := &RateLimiterConfig{
		Token: &RateLimiterRateConfig{
			MaxRequestsPerSecond:  10,
			BlockTimeMilliseconds: 100,
		},
		CustomTokens: &map[string]*RateLimiterRateConfig{},
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("DONE"))
	})

	rateLimiterCheckFunction := func(ctx context.Context, keyType string, key string, config *RateLimiterConfig, rateConfig *RateLimiterRateConfig) (*time.Time, error) {
		return nil, nil
	}

	request := httptest.NewRequest("GET", "http://testing", nil)
	request.Header.Add("API_KEY", "123")
	recorder := httptest.NewRecorder()

	rateLimiter(config, nextHandler, rateLimiterCheckFunction).ServeHTTP(recorder, request)

	response := recorder.Result()
	responseStatus := response.StatusCode
	responseBody, err := ioutil.ReadAll(response.Body)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 200, responseStatus)
	assert.Equal(s.T(), "DONE", string(responseBody))
}
