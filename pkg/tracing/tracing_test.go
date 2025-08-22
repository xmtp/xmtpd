package tracing

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap/zaptest"
)

func Test_GoPanicWrap_WaitGroup(t *testing.T) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	finished := false
	var finishedLock sync.RWMutex
	GoPanicWrap(ctx, &wg, "test", func(ctx context.Context) {
		<-ctx.Done()
		finishedLock.Lock()
		defer finishedLock.Unlock()
		finished = true
	})
	done := false
	var doneLock sync.RWMutex
	go func() {
		wg.Wait()
		doneLock.Lock()
		defer doneLock.Unlock()
		done = true
	}()
	go func() { time.Sleep(time.Millisecond); cancel() }()

	assert.Eventually(t, func() bool {
		finishedLock.RLock()
		defer finishedLock.RUnlock()
		doneLock.RLock()
		defer doneLock.RUnlock()
		return finished && done
	}, time.Second, 10*time.Millisecond)
}

func TestOTelConfig_Integration(t *testing.T) {
	config := EnableDevTracer()
	config.UseStdout = false

	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	cleanup, err := InitializeOTel(ctx, config, logger)
	require.NoError(t, err)
	defer cleanup()

	// Override with in-memory tracer for verification
	tp, getSpans := NewInMemoryTracer()
	defer func() {
		_ = tp.Shutdown(context.Background())
	}()

	_, span := StartOTelSpan(ctx, "integration.test")
	SetAttributes(ctx, attribute.String("test.type", "integration"))
	span.End()

	spans := getSpans()
	assert.Len(t, spans, 1)
	assert.Equal(t, "integration.test", spans[0].Name())
}

func TestParseResourceAttributes(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected []attribute.KeyValue
	}{
		{
			name:     "empty environment variable",
			envValue: "",
			expected: nil,
		},
		{
			name:     "single attribute",
			envValue: "environment=development",
			expected: []attribute.KeyValue{
				attribute.String("environment", "development"),
			},
		},
		{
			name:     "multiple attributes",
			envValue: "environment=production,version=1.0.0,region=us-west-2",
			expected: []attribute.KeyValue{
				attribute.String("environment", "production"),
				attribute.String("version", "1.0.0"),
				attribute.String("region", "us-west-2"),
			},
		},
		{
			name:     "attributes with spaces",
			envValue: " environment = development , version = 1.0.0 ",
			expected: []attribute.KeyValue{
				attribute.String("environment", "development"),
				attribute.String("version", "1.0.0"),
			},
		},
		{
			name:     "value with equals sign",
			envValue: "sql_query=SELECT * FROM users WHERE id = 123",
			expected: []attribute.KeyValue{
				attribute.String("sql_query", "SELECT * FROM users WHERE id = 123"),
			},
		},
		{
			name:     "invalid format ignored",
			envValue: "valid=value,invalid_no_equals,another=valid",
			expected: []attribute.KeyValue{
				attribute.String("valid", "value"),
				attribute.String("another", "valid"),
			},
		},
		{
			name:     "empty keys and values ignored",
			envValue: "=empty_key,empty_value=,valid=value",
			expected: []attribute.KeyValue{
				attribute.String("valid", "value"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variable
			originalValue := os.Getenv("OTEL_RESOURCE_ATTRIBUTES")
			if tt.envValue != "" {
				t.Setenv("OTEL_RESOURCE_ATTRIBUTES", tt.envValue)
			} else {
				_ = os.Unsetenv("OTEL_RESOURCE_ATTRIBUTES")
			}
			defer func() {
				if originalValue != "" {
					_ = os.Setenv("OTEL_RESOURCE_ATTRIBUTES", originalValue)
				} else {
					_ = os.Unsetenv("OTEL_RESOURCE_ATTRIBUTES")
				}
			}()

			// Test the parsing function
			result := parseResourceAttributes()

			// Verify results
			assert.Equal(t, len(tt.expected), len(result), "Number of attributes should match")
			for i, expected := range tt.expected {
				if i < len(result) {
					assert.Equal(t, expected.Key, result[i].Key, "Attribute key should match")
					assert.Equal(t,
						expected.Value.AsString(),
						result[i].Value.AsString(),
						"Attribute value should match")
				}
			}
		})
	}
}
