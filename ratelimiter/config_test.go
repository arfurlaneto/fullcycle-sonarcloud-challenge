package ratelimiter

import (
	"os"
	"testing"

	"github.com/arfurlaneto/fullcycle-sonarcloud-challenge/ratelimiter/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ConfigTestSuite struct {
	suite.Suite
	controller *gomock.Controller
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) SetupTest() {
	s.controller = gomock.NewController(s.T())
	os.Unsetenv(envKeyIPMaxRequestsPerSecond)
	os.Unsetenv(envKeyIPBlockTimeMilliseconds)
	os.Unsetenv(envKeyTokenMaxRequestsPerSecond)
	os.Unsetenv(envKeyTokenBlockTimeMilliseconds)
	os.Unsetenv(envKeyDebug)
	os.Unsetenv(envUseRedis)
	os.Unsetenv(envRedisAddress)
	os.Unsetenv(envRedisPassword)
	os.Unsetenv(envRedisDB)
	os.Unsetenv("RATE_LIMITER_TOKEN_abc_MAX_REQUESTS")
	os.Unsetenv("RATE_LIMITER_TOKEN_abc_BLOCK_TIME")
	os.Unsetenv("RATE_LIMITER_TOKEN_def_MAX_REQUESTS")
	os.Unsetenv("RATE_LIMITER_TOKEN_def_BLOCK_TIME")
}

func (s *ConfigTestSuite) TestGetDefaultConfiguration() {
	config := getDefaultConfiguration()
	assert.NotNil(s.T(), config)
	assert.NotNil(s.T(), config.IP)
	assert.NotNil(s.T(), config.Token)
	assert.NotNil(s.T(), config.CustomTokens)
	assert.Empty(s.T(), config.CustomTokens)
	assert.NotNil(s.T(), config.StorageAdapter)
	assert.NotNil(s.T(), config.ResponseWriter)
	assert.Equal(s.T(), false, config.Debug)
}

func (s *ConfigTestSuite) TestSetConfiguration_AllSeted() {

	storageAdapterMock := mocks.NewMockRateLimitStorageAdapter(s.controller)
	responseWriterMock := mocks.NewMockRateLimiterResponseWriter(s.controller)

	inputConfig := &RateLimiterConfig{
		IP: &RateLimiterRateConfig{
			MaxRequestsPerSecond:  111,
			BlockTimeMilliseconds: 222,
		},
		Token: &RateLimiterRateConfig{
			MaxRequestsPerSecond:  333,
			BlockTimeMilliseconds: 444,
		},
		CustomTokens: &map[string]*RateLimiterRateConfig{
			"abc": {MaxRequestsPerSecond: 555, BlockTimeMilliseconds: 666},
			"def": {MaxRequestsPerSecond: 777, BlockTimeMilliseconds: 888},
		},
		StorageAdapter: storageAdapterMock,
		ResponseWriter: responseWriterMock,
		Debug:          true,
	}

	config := setConfiguration(inputConfig)
	assert.NotNil(s.T(), config)
	assert.Equal(s.T(), int64(111), config.IP.MaxRequestsPerSecond)
	assert.Equal(s.T(), int64(222), config.IP.BlockTimeMilliseconds)
	assert.Equal(s.T(), int64(333), config.Token.MaxRequestsPerSecond)
	assert.Equal(s.T(), int64(444), config.Token.BlockTimeMilliseconds)
	assert.NotNil(s.T(), config.CustomTokens)
	assert.Len(s.T(), *config.CustomTokens, 2)
	assert.Contains(s.T(), *config.CustomTokens, "abc")
	assert.Contains(s.T(), *config.CustomTokens, "def")
	assert.Equal(s.T(), (*config.CustomTokens)["abc"].MaxRequestsPerSecond, int64(555))
	assert.Equal(s.T(), (*config.CustomTokens)["abc"].BlockTimeMilliseconds, int64(666))
	assert.Equal(s.T(), (*config.CustomTokens)["def"].MaxRequestsPerSecond, int64(777))
	assert.Equal(s.T(), (*config.CustomTokens)["def"].BlockTimeMilliseconds, int64(888))
	assert.Equal(s.T(), storageAdapterMock, config.StorageAdapter)
	assert.Equal(s.T(), responseWriterMock, config.ResponseWriter)
	assert.Equal(s.T(), true, config.Debug)
}

func (s *ConfigTestSuite) TestSetConfiguration_AllSeted_EmptyCustomTokenConfigs() {

	storageAdapterMock := mocks.NewMockRateLimitStorageAdapter(s.controller)
	responseWriterMock := mocks.NewMockRateLimiterResponseWriter(s.controller)

	inputConfig := &RateLimiterConfig{
		IP: &RateLimiterRateConfig{
			MaxRequestsPerSecond:  111,
			BlockTimeMilliseconds: 222,
		},
		Token: &RateLimiterRateConfig{
			MaxRequestsPerSecond:  333,
			BlockTimeMilliseconds: 444,
		},
		CustomTokens: &map[string]*RateLimiterRateConfig{
			"abc": nil,
			"def": nil,
		},
		StorageAdapter: storageAdapterMock,
		ResponseWriter: responseWriterMock,
		Debug:          true,
	}

	config := setConfiguration(inputConfig)
	assert.NotNil(s.T(), config)
	assert.Equal(s.T(), int64(111), config.IP.MaxRequestsPerSecond)
	assert.Equal(s.T(), int64(222), config.IP.BlockTimeMilliseconds)
	assert.Equal(s.T(), int64(333), config.Token.MaxRequestsPerSecond)
	assert.Equal(s.T(), int64(444), config.Token.BlockTimeMilliseconds)
	assert.NotNil(s.T(), config.CustomTokens)
	assert.Len(s.T(), *config.CustomTokens, 2)
	assert.Contains(s.T(), *config.CustomTokens, "abc")
	assert.Contains(s.T(), *config.CustomTokens, "def")
	assert.Equal(s.T(), (*config.CustomTokens)["abc"].MaxRequestsPerSecond, int64(333))
	assert.Equal(s.T(), (*config.CustomTokens)["abc"].BlockTimeMilliseconds, int64(444))
	assert.Equal(s.T(), (*config.CustomTokens)["def"].MaxRequestsPerSecond, int64(333))
	assert.Equal(s.T(), (*config.CustomTokens)["def"].BlockTimeMilliseconds, int64(444))
	assert.Equal(s.T(), storageAdapterMock, config.StorageAdapter)
	assert.Equal(s.T(), responseWriterMock, config.ResponseWriter)
	assert.Equal(s.T(), true, config.Debug)
}

func (s *ConfigTestSuite) TestSetConfiguration_NilInput() {
	config := setConfiguration(nil)
	assert.NotNil(s.T(), config)
	assert.NotNil(s.T(), config.IP)
	assert.NotNil(s.T(), config.Token)
	assert.NotNil(s.T(), config.CustomTokens)
	assert.Empty(s.T(), config.CustomTokens)
	assert.NotNil(s.T(), config.StorageAdapter)
	assert.NotNil(s.T(), config.ResponseWriter)
	assert.Equal(s.T(), false, config.Debug)
}

func (s *ConfigTestSuite) TestSetConfiguration_EmptyInput() {
	config := setConfiguration(&RateLimiterConfig{})
	assert.NotNil(s.T(), config)
	assert.NotNil(s.T(), config.IP)
	assert.NotNil(s.T(), config.Token)
	assert.NotNil(s.T(), config.CustomTokens)
	assert.Empty(s.T(), config.CustomTokens)
	assert.NotNil(s.T(), config.StorageAdapter)
	assert.NotNil(s.T(), config.ResponseWriter)
	assert.Equal(s.T(), false, config.Debug)
}

func (s *ConfigTestSuite) TestSetConfiguration_ValuesFromEnv() {
	os.Setenv(envKeyIPMaxRequestsPerSecond, "111")
	os.Setenv(envKeyIPBlockTimeMilliseconds, "222")
	os.Setenv(envKeyTokenMaxRequestsPerSecond, "333")
	os.Setenv(envKeyTokenBlockTimeMilliseconds, "444")
	os.Setenv("RATE_LIMITER_TOKEN_abc_MAX_REQUESTS", "555")
	os.Setenv("RATE_LIMITER_TOKEN_abc_BLOCK_TIME", "666")
	os.Setenv("RATE_LIMITER_TOKEN_def_MAX_REQUESTS", "777")
	os.Setenv("RATE_LIMITER_TOKEN_def_BLOCK_TIME", "888")
	os.Setenv(envKeyDebug, "true")

	config := setConfiguration(nil)
	assert.NotNil(s.T(), config)
	assert.Equal(s.T(), int64(111), config.IP.MaxRequestsPerSecond)
	assert.Equal(s.T(), int64(222), config.IP.BlockTimeMilliseconds)
	assert.Equal(s.T(), int64(333), config.Token.MaxRequestsPerSecond)
	assert.Equal(s.T(), int64(444), config.Token.BlockTimeMilliseconds)
	assert.NotNil(s.T(), config.CustomTokens)
	assert.Len(s.T(), *config.CustomTokens, 2)
	assert.Contains(s.T(), *config.CustomTokens, "abc")
	assert.Contains(s.T(), *config.CustomTokens, "def")
	assert.Equal(s.T(), (*config.CustomTokens)["abc"].MaxRequestsPerSecond, int64(555))
	assert.Equal(s.T(), (*config.CustomTokens)["abc"].BlockTimeMilliseconds, int64(666))
	assert.Equal(s.T(), (*config.CustomTokens)["def"].MaxRequestsPerSecond, int64(777))
	assert.Equal(s.T(), (*config.CustomTokens)["def"].BlockTimeMilliseconds, int64(888))
	assert.NotNil(s.T(), config.StorageAdapter)
	assert.NotNil(s.T(), config.ResponseWriter)
	assert.Equal(s.T(), true, config.Debug)
}

func (s *ConfigTestSuite) TestSetConfiguration_ValuesFromEnv_CustomTokenDefaultsToToken() {

	os.Setenv(envKeyTokenMaxRequestsPerSecond, "333")
	os.Setenv(envKeyTokenBlockTimeMilliseconds, "444")
	os.Setenv("RATE_LIMITER_TOKEN_abc_MAX_REQUESTS", "555")
	os.Setenv("RATE_LIMITER_TOKEN_def_BLOCK_TIME", "888")

	config := setConfiguration(nil)
	assert.NotNil(s.T(), config)

	assert.NotNil(s.T(), config.CustomTokens)
	assert.Len(s.T(), *config.CustomTokens, 2)
	assert.Contains(s.T(), *config.CustomTokens, "abc")
	assert.Contains(s.T(), *config.CustomTokens, "def")
	assert.Equal(s.T(), (*config.CustomTokens)["abc"].MaxRequestsPerSecond, int64(555))
	assert.Equal(s.T(), (*config.CustomTokens)["abc"].BlockTimeMilliseconds, int64(444))
	assert.Equal(s.T(), (*config.CustomTokens)["def"].MaxRequestsPerSecond, int64(333))
	assert.Equal(s.T(), (*config.CustomTokens)["def"].BlockTimeMilliseconds, int64(888))
}

func (s *ConfigTestSuite) TestSetConfiguration_RedisAdapter() {
	os.Setenv(envUseRedis, "true")
	os.Setenv(envRedisAddress, "localhost:6379")
	os.Setenv(envRedisPassword, "")
	os.Setenv(envRedisDB, "")

	config := setConfiguration(nil)
	assert.NotNil(s.T(), config)
	assert.NotNil(s.T(), config.IP.MaxRequestsPerSecond)
	assert.NotNil(s.T(), config.IP.BlockTimeMilliseconds)
	assert.NotNil(s.T(), config.Token.MaxRequestsPerSecond)
	assert.NotNil(s.T(), config.Token.BlockTimeMilliseconds)
	assert.NotNil(s.T(), config.StorageAdapter)
	assert.NotNil(s.T(), config.ResponseWriter)
	assert.Equal(s.T(), false, config.Debug)
}

func (s *ConfigTestSuite) TestSetConfiguration_RedisAdapterErrMissingAddress() {
	os.Setenv(envUseRedis, "true")
	assert.Panics(s.T(), func() { setConfiguration(nil) }, "should panic")
}

func (s *ConfigTestSuite) TestGetRateLimiterRateConfigForToken() {
	inputConfig := &RateLimiterConfig{
		Token: &RateLimiterRateConfig{
			MaxRequestsPerSecond:  333,
			BlockTimeMilliseconds: 444,
		},
		CustomTokens: &map[string]*RateLimiterRateConfig{
			"abc": {MaxRequestsPerSecond: 555, BlockTimeMilliseconds: 666},
		},
	}

	config := setConfiguration(inputConfig)

	abcConfig, abcIsCustom := config.GetRateLimiterRateConfigForToken("abc")
	zzzConfig, zzzIsCustom := config.GetRateLimiterRateConfigForToken("zzz")

	assert.Equal(s.T(), int64(555), abcConfig.MaxRequestsPerSecond)
	assert.Equal(s.T(), int64(666), abcConfig.BlockTimeMilliseconds)
	assert.Equal(s.T(), true, abcIsCustom)
	assert.Equal(s.T(), int64(333), zzzConfig.MaxRequestsPerSecond)
	assert.Equal(s.T(), int64(444), zzzConfig.BlockTimeMilliseconds)
	assert.Equal(s.T(), false, zzzIsCustom)
}
