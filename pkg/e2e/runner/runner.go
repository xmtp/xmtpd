package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/xmtp/xmtpd/pkg/e2e/types"
	"go.uber.org/zap"
)

type Config = types.Config

type TestInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TestResult struct {
	Name     string        `json:"name"`
	Status   string        `json:"status"`
	Duration time.Duration `json:"duration"`
	Error    string        `json:"error,omitempty"`
}

type Runner struct {
	logger *zap.Logger
	cfg    Config
	tests  []Test
}

func New(logger *zap.Logger, cfg Config) *Runner {
	return &Runner{
		logger: logger,
		cfg:    cfg,
		tests:  AllTests(),
	}
}

func (r *Runner) Run(ctx context.Context) error {
	selected := r.filterTests()
	if len(selected) == 0 {
		r.logger.Warn("no tests matched the filter", zap.Strings("filter", r.cfg.TestFilter))
		return errors.New("no tests matched")
	}

	r.logger.Info("starting e2e test run",
		zap.Int("total_tests", len(selected)),
		zap.String("xmtpd_image", r.cfg.XmtpdImage),
	)

	results := make([]TestResult, 0, len(selected))
	var failures int

	for _, t := range selected {
		result := r.runTest(ctx, t)
		results = append(results, result)
		if result.Status == "FAIL" {
			failures++
		}
	}

	r.printResults(results)

	if failures > 0 {
		return fmt.Errorf("%d test(s) failed", failures)
	}
	return nil
}

func (r *Runner) runTest(ctx context.Context, t Test) TestResult {
	logger := r.logger.Named(t.Name())
	logger.Info("running test")

	start := time.Now()

	env, err := NewEnvironment(ctx, logger, r.cfg)
	if err != nil {
		duration := time.Since(start)
		logger.Error("failed to set up environment", zap.Error(err))
		return TestResult{
			Name:     t.Name(),
			Status:   "FAIL",
			Duration: duration,
			Error:    fmt.Sprintf("environment setup: %s", err),
		}
	}

	// Inject TestingT so tests can use require/assert
	testingT := types.NewTestingT(logger)
	env.SetTestingT(testingT)

	testErr := r.executeTest(ctx, t, env)
	duration := time.Since(start)

	cleanupErr := env.Cleanup(ctx)
	if cleanupErr != nil {
		logger.Warn("environment cleanup error", zap.Error(cleanupErr))
	}

	if testErr != nil {
		logger.Error("test failed", zap.Error(testErr), zap.Duration("duration", duration))
		return TestResult{
			Name:     t.Name(),
			Status:   "FAIL",
			Duration: duration,
			Error:    testErr.Error(),
		}
	}

	logger.Info("test passed", zap.Duration("duration", duration))
	return TestResult{
		Name:     t.Name(),
		Status:   "PASS",
		Duration: duration,
	}
}

// executeTest runs the test and recovers from TestFailedError panics caused by
// require.FailNow(). This allows tests to use testify's require package which
// calls FailNow() on assertion failure.
func (r *Runner) executeTest(ctx context.Context, t Test, env *Environment) (testErr error) {
	defer func() {
		if r := recover(); r != nil {
			if tfe, ok := r.(*types.TestFailedError); ok {
				testErr = tfe
			} else {
				// Re-panic for unexpected panics (bugs, nil derefs, etc.)
				panic(r)
			}
		}
	}()
	return t.Run(ctx, env)
}

func (r *Runner) filterTests() []Test {
	if len(r.cfg.TestFilter) == 0 {
		return r.tests
	}
	filterSet := make(map[string]struct{}, len(r.cfg.TestFilter))
	for _, f := range r.cfg.TestFilter {
		filterSet[f] = struct{}{}
	}
	var selected []Test
	for _, t := range r.tests {
		if _, ok := filterSet[t.Name()]; ok {
			selected = append(selected, t)
		}
	}
	return selected
}

func (r *Runner) printResults(results []TestResult) {
	if r.cfg.OutputFormat == "json" {
		r.printJSON(results)
		return
	}
	r.printText(results)
}

func (r *Runner) printText(results []TestResult) {
	fmt.Println("\n=== E2E Test Results ===")
	for _, res := range results {
		status := "PASS"
		if res.Status == "FAIL" {
			status = "FAIL"
		}
		fmt.Printf("  [%s] %s (%s)\n", status, res.Name, res.Duration.Truncate(time.Millisecond))
		if res.Error != "" {
			fmt.Printf("         error: %s\n", res.Error)
		}
	}
	fmt.Println()
}

func (r *Runner) printJSON(results []TestResult) {
	data, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(data))
}

func (r *Runner) Tests() []TestInfo {
	infos := make([]TestInfo, len(r.tests))
	for i, t := range r.tests {
		infos[i] = TestInfo{
			Name:        t.Name(),
			Description: t.Description(),
		}
	}
	return infos
}
