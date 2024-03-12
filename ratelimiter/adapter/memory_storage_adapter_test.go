package adapter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RateLimitMemoryStorageAdapterTestSuite struct {
	suite.Suite
	context context.Context
}

func TestRateLimitMemoryStorageAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(RateLimitMemoryStorageAdapterTestSuite))
}

func (s *RateLimitMemoryStorageAdapterTestSuite) SetupTest() {
	s.context = context.Background()
}

func (s *RateLimitMemoryStorageAdapterTestSuite) TestNewRateLimitMemoryStorageAdapter() {
	storageAdapter := NewRateLimitMemoryStorageAdapter()
	assert.NotNil(s.T(), storageAdapter)
}

func (s *RateLimitMemoryStorageAdapterTestSuite) TestIncrementAccesses() {
	ctx := s.context
	keyType := "IP"
	keyValue := "127.0.0.1"
	maxAccesses := 5

	storageAdapter := NewRateLimitMemoryStorageAdapter()

	expectedResults := [][]interface{}{
		{true, int64(1), nil},
		{true, int64(2), nil},
		{true, int64(3), nil},
		{true, int64(4), nil},
		{true, int64(5), nil},
		{false, int64(5), nil},
	}

	for _, val := range expectedResults {
		success, count, err := storageAdapter.IncrementAccesses(ctx, keyType, keyValue, int64(maxAccesses))
		assert.Equal(s.T(), val[0], success)
		assert.Equal(s.T(), val[1], count)
		assert.Equal(s.T(), val[2], err)
	}
}

func (s *RateLimitMemoryStorageAdapterTestSuite) TestAddBlockGetBlock_SameTypeAndValue() {
	ctx := s.context
	keyType := "IP"
	keyValue := "127.0.0.1"

	storageAdapter := NewRateLimitMemoryStorageAdapter()

	addBlockResult, addBlockErr := storageAdapter.AddBlock(ctx, keyType, keyValue, 100)
	getBlockResult, getBlockErr := storageAdapter.GetBlock(ctx, keyType, keyValue)

	assert.Nil(s.T(), addBlockErr)
	assert.Nil(s.T(), getBlockErr)
	assert.NotNil(s.T(), addBlockResult)
	assert.NotNil(s.T(), getBlockResult)
	assert.Equal(s.T(), addBlockResult, getBlockResult)
}

func (s *RateLimitMemoryStorageAdapterTestSuite) TestAddBlockGetBlock_AnotherType() {
	ctx := s.context
	keyType := "IP"
	keyValue := "127.0.0.1"

	storageAdapter := NewRateLimitMemoryStorageAdapter()

	addBlockResult, addBlockErr := storageAdapter.AddBlock(ctx, keyType, keyValue, 100)
	getBlockResult, getBlockErr := storageAdapter.GetBlock(ctx, "TOKEN", keyValue)

	assert.Nil(s.T(), addBlockErr)
	assert.Nil(s.T(), getBlockErr)
	assert.NotNil(s.T(), addBlockResult)
	assert.Nil(s.T(), getBlockResult)
	assert.NotEqual(s.T(), addBlockResult, getBlockResult)
}

func (s *RateLimitMemoryStorageAdapterTestSuite) TestAddBlockGetBlock_AnotherValue() {
	ctx := s.context
	keyType := "IP"
	keyValue := "127.0.0.1"

	storageAdapter := NewRateLimitMemoryStorageAdapter()

	addBlockResult, addBlockErr := storageAdapter.AddBlock(ctx, keyType, keyValue, 100)
	getBlockResult, getBlockErr := storageAdapter.GetBlock(ctx, keyType, "127.0.0.2")

	assert.Nil(s.T(), addBlockErr)
	assert.Nil(s.T(), getBlockErr)
	assert.NotNil(s.T(), addBlockResult)
	assert.Nil(s.T(), getBlockResult)
	assert.NotEqual(s.T(), addBlockResult, getBlockResult)
}

func (s *RateLimitMemoryStorageAdapterTestSuite) TestAddBlockGetBlock_ExpiredBlock() {
	ctx := s.context
	keyType := "IP"
	keyValue := "127.0.0.1"

	storageAdapter := NewRateLimitMemoryStorageAdapter()

	addBlockResult, addBlockErr := storageAdapter.AddBlock(ctx, keyType, keyValue, -100)
	getBlockResult, getBlockErr := storageAdapter.GetBlock(ctx, keyType, keyValue)

	assert.Nil(s.T(), addBlockErr)
	assert.Nil(s.T(), getBlockErr)
	assert.NotNil(s.T(), addBlockResult)
	assert.Nil(s.T(), getBlockResult)
	assert.NotEqual(s.T(), addBlockResult, getBlockResult)
}
