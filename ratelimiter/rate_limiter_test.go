package ratelimiter

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/arfurlaneto/fullcycle-sonarcloud-challenge/ratelimiter/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type RateLimiterTestSuite struct {
	suite.Suite
	controller         *gomock.Controller
	context            context.Context
	storageAdapterMock *mocks.MockRateLimitStorageAdapter
}

func TestRateLimiterTestSuite(t *testing.T) {
	suite.Run(t, new(RateLimiterTestSuite))
}

func (s *RateLimiterTestSuite) SetupTest() {
	s.controller = gomock.NewController(s.T())
	s.context = context.Background()
	s.storageAdapterMock = mocks.NewMockRateLimitStorageAdapter(s.controller)
}

func (s *RateLimiterTestSuite) TestCheckRateLimit_AccessAllowed() {
	context := s.context
	keyType := "IP"
	key := "127.0.0.1"
	config := &RateLimiterConfig{
		IP: &RateLimiterRateConfig{
			MaxRequestsPerSecond:  10,
			BlockTimeMilliseconds: 100,
		},
	}

	s.storageAdapterMock.EXPECT().
		GetBlock(context, keyType, key).Return(nil, nil).Times(1)

	s.storageAdapterMock.EXPECT().
		IncrementAccesses(context, keyType, key, gomock.Any()).Return(true, int64(1), nil).Times(1)

	config.StorageAdapter = s.storageAdapterMock

	returnedBlock, err := checkRateLimit(context, keyType, key, config, config.IP)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), returnedBlock)
}

func (s *RateLimiterTestSuite) TestCheckRateLimit_AccessDenied() {
	context := s.context
	keyType := "IP"
	key := "127.0.0.1"
	config := &RateLimiterConfig{
		IP: &RateLimiterRateConfig{
			MaxRequestsPerSecond:  10,
			BlockTimeMilliseconds: 100,
		},
	}
	block := time.Now().Add(time.Millisecond * 100)

	s.storageAdapterMock.EXPECT().
		GetBlock(context, keyType, key).Return(nil, nil).Times(1)

	s.storageAdapterMock.EXPECT().
		IncrementAccesses(context, keyType, key, gomock.Any()).Return(false, int64(10), nil).Times(1)

	s.storageAdapterMock.EXPECT().
		AddBlock(context, keyType, key, config.IP.BlockTimeMilliseconds).Return(&block, nil).Times(1)

	config.StorageAdapter = s.storageAdapterMock

	returnedBlock, err := checkRateLimit(context, keyType, key, config, config.IP)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), block, *returnedBlock)
}

func (s *RateLimiterTestSuite) TestCheckRateLimit_AlreadyBlocked() {
	context := s.context
	keyType := "IP"
	key := "127.0.0.1"
	config := &RateLimiterConfig{
		IP: &RateLimiterRateConfig{
			MaxRequestsPerSecond:  10,
			BlockTimeMilliseconds: 100,
		},
	}
	block := time.Now().Add(time.Millisecond * 100)

	s.storageAdapterMock.EXPECT().
		GetBlock(context, keyType, key).Return(&block, nil).Times(1)

	config.StorageAdapter = s.storageAdapterMock

	returnedBlock, err := checkRateLimit(context, keyType, key, config, config.IP)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), block, *returnedBlock)
}

func (s *RateLimiterTestSuite) TestCheckRateLimit_EmptyKey() {
	context := s.context
	keyType := "IP"
	key := ""
	config := &RateLimiterConfig{
		IP: &RateLimiterRateConfig{
			MaxRequestsPerSecond:  10,
			BlockTimeMilliseconds: 100,
		},
	}

	config.StorageAdapter = s.storageAdapterMock

	returnedBlock, err := checkRateLimit(context, keyType, key, config, config.IP)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), returnedBlock)
}

func (s *RateLimiterTestSuite) TestCheckRateLimit_GetBlockError() {
	context := s.context
	keyType := "IP"
	key := "127.0.0.1"
	config := &RateLimiterConfig{
		IP: &RateLimiterRateConfig{
			MaxRequestsPerSecond:  10,
			BlockTimeMilliseconds: 100,
		},
	}

	s.storageAdapterMock.EXPECT().
		GetBlock(context, keyType, key).Return(nil, errors.New("error")).Times(1)

	config.StorageAdapter = s.storageAdapterMock

	returnedBlock, err := checkRateLimit(context, keyType, key, config, config.IP)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), returnedBlock)
}

func (s *RateLimiterTestSuite) TestCheckRateLimit_IncrementAccessesError() {
	context := s.context
	keyType := "IP"
	key := "127.0.0.1"
	config := &RateLimiterConfig{
		IP: &RateLimiterRateConfig{
			MaxRequestsPerSecond:  10,
			BlockTimeMilliseconds: 100,
		},
	}

	s.storageAdapterMock.EXPECT().
		GetBlock(context, keyType, key).Return(nil, nil).Times(1)

	s.storageAdapterMock.EXPECT().
		IncrementAccesses(context, keyType, key, gomock.Any()).Return(false, int64(1), errors.New("error")).Times(1)

	config.StorageAdapter = s.storageAdapterMock

	returnedBlock, err := checkRateLimit(context, keyType, key, config, config.IP)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), returnedBlock)
}

func (s *RateLimiterTestSuite) TestCheckRateLimit_AddBlockError() {
	context := s.context
	keyType := "IP"
	key := "127.0.0.1"
	config := &RateLimiterConfig{
		IP: &RateLimiterRateConfig{
			MaxRequestsPerSecond:  10,
			BlockTimeMilliseconds: 100,
		},
	}

	s.storageAdapterMock.EXPECT().
		GetBlock(context, keyType, key).Return(nil, nil).Times(1)

	s.storageAdapterMock.EXPECT().
		IncrementAccesses(context, keyType, key, gomock.Any()).Return(false, int64(10), nil).Times(1)

	s.storageAdapterMock.EXPECT().
		AddBlock(context, keyType, key, config.IP.BlockTimeMilliseconds).Return(nil, errors.New("error")).Times(1)

	config.StorageAdapter = s.storageAdapterMock

	returnedBlock, err := checkRateLimit(context, keyType, key, config, config.IP)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), returnedBlock)
}
