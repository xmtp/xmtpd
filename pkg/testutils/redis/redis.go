// Package redis implements the redis test utils.
package redis

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

const redisAddress = "localhost:6379"

// NewRedisForTest creates a Redis client configured for testing with proper cleanup.
// It automatically generates a unique key prefix based on the test name to avoid conflicts.
// All keys created with this prefix are cleaned up after the test.
func NewRedisForTest(t *testing.T) (redis.UniversalClient, string) {
	ctx := context.Background()

	keyPrefix := generateTestKeyPrefix(t)

	// Create Redis client with default test configuration
	client := redis.NewClient(&redis.Options{
		Addr: redisAddress,
		DB:   15, // Use DB 15 for tests
	})

	// Test connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		_ = client.Close()
		t.Errorf("Redis not available: %v", err)
		return nil, ""
	}

	// Setup cleanup
	t.Cleanup(func() {
		cleanupKeysByPrefix(t, client, keyPrefix)
		_ = client.Close()
	})

	t.Logf("Redis test client created with key prefix: %s", keyPrefix)
	return client, keyPrefix
}

// generateTestKeyPrefix creates a unique key prefix based on test name and timestamp
func generateTestKeyPrefix(t *testing.T) string {
	// Clean test name to be Redis-key safe
	testName := strings.ReplaceAll(t.Name(), "/", "_")
	testName = strings.ReplaceAll(testName, " ", "_")
	testName = strings.ReplaceAll(testName, "-", "_")

	// Add timestamp to ensure uniqueness even for parallel runs
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	return fmt.Sprintf("test:%s:%d:", testName, timestamp)
}

// cleanupKeysByPrefix removes all keys matching the prefix pattern
func cleanupKeysByPrefix(t *testing.T, client redis.UniversalClient, keyPrefix string) {
	ctx := context.Background()
	pattern := keyPrefix + "*"

	// Use SCAN to find keys matching the pattern
	iter := client.Scan(ctx, 0, pattern, 0).Iterator()
	var keysToDelete []string

	for iter.Next(ctx) {
		keysToDelete = append(keysToDelete, iter.Val())
	}

	if err := iter.Err(); err != nil {
		t.Logf("Warning: Error scanning Redis keys with pattern %s: %v", pattern, err)
		return
	}

	// Delete found keys
	if len(keysToDelete) > 0 {
		deletedCount, err := client.Del(ctx, keysToDelete...).Result()
		if err != nil {
			t.Logf("Warning: Failed to cleanup Redis keys with pattern %s: %v", pattern, err)
		} else if deletedCount > 0 {
			t.Logf("Cleaned up %d Redis keys with prefix: %s", deletedCount, keyPrefix)
		}
	}
}

// FlushRedisForTest safely flushes a Redis database for testing
// This should only be used when you need to completely reset the test database
func FlushRedisForTest(t *testing.T, client redis.UniversalClient) {
	ctx := context.Background()

	err := client.FlushDB(ctx).Err()
	if err != nil {
		t.Fatalf("Failed to flush Redis database for test: %v", err)
	}

	t.Log("Flushed Redis database for test")
}
