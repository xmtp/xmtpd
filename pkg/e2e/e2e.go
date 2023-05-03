package e2e

import (
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/zap"
	"golang.org/x/sync/errgroup"

	_ "net/http/pprof"
)

type E2E struct {
	ctx     context.Context
	log     *zap.Logger
	metrics *Metrics
	rand    *rand.Rand
	randMu  sync.Mutex

	opts *Options
}

type Options struct {
	APIURLs          []string      `long:"api-url" env:"XMTP_API_URLS" description:"XMTP node API URLs" default:"http://localhost"`
	ClientsPerURL    int           `long:"clients-per-url" description:"Number of clients for each API URL" default:"1"`
	MessagePerClient int           `long:"messages-per-client" description:"Number of messages to publish for each client" default:"3"`
	Continuous       bool          `long:"continuous" description:"Run continuously"`
	ExitOnError      bool          `long:"exit-on-error" description:"Exit on error if running continuously"`
	RunDelay         time.Duration `long:"delay" description:"Delay between runs (in seconds)" default:"5s"`
	AdminPort        uint          `long:"admin-port" description:"Admin HTTP server listen port" default:"7777"`

	GitCommit string
}

type testRunFunc func(name string) error

type Test struct {
	Name string
	Run  testRunFunc
}

func (e *E2E) Tests() []*Test {
	return []*Test{
		e.newTest("convergence", e.testConvergence),
	}
}

func New(ctx context.Context, opts *Options) (*E2E, error) {
	e := &E2E{
		ctx:     ctx,
		log:     ctx.Logger().Named("e2e"),
		metrics: newMetrics(),
		rand:    rand.New(rand.NewSource(time.Now().UTC().UnixNano())),
		opts:    opts,
	}
	e.log.Info("running", zap.String("git-commit", opts.GitCommit), zap.Strings("nodes", opts.APIURLs))

	if e.opts.Continuous {
		go func() {
			// Initialize HTTP server for profiler and metrics.
			http.Handle("/metrics", promhttp.Handler())
			http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {})
			addr := net.JoinHostPort("0.0.0.0", strconv.Itoa(int(opts.AdminPort)))
			go func() {
				e.log.Info("serving admin http", zap.String("address", addr))
				err := http.ListenAndServe(addr, nil)
				if err != nil {
					e.log.Error("serving e2e admin", zap.Error(err))
				}
			}()
		}()
	}

	for {
		g, _ := errgroup.WithContext(e.ctx)
		for _, test := range e.Tests() {
			test := test
			g.Go(func() error {
				return e.runTest(test)
			})
		}
		err := g.Wait()
		if err != nil && (!e.opts.Continuous || e.opts.ExitOnError) {
			return nil, err
		}
		if !e.opts.Continuous {
			break
		}
		time.Sleep(e.opts.RunDelay)
	}

	return e, nil
}

func (e *E2E) runTest(test *Test) error {
	started := time.Now().UTC()
	log := e.log.With(zap.String("test", test.Name))

	err := test.Run(test.Name)
	duration := time.Since(started)
	log = log.With(zap.Duration("duration", duration))
	if err != nil {
		log.Error("test failed", zap.Error(err))
		e.metrics.recordRun(e.ctx, test.Name, "failed", duration)
		return err
	}
	log.Info("test passed")
	e.metrics.recordRun(e.ctx, test.Name, "passed", duration)

	return nil
}

func (e *E2E) newTest(name string, runFn testRunFunc) *Test {
	return &Test{
		Name: name,
		Run:  runFn,
	}
}

func (e *E2E) randomStringLower(n int) string {
	e.randMu.Lock()
	defer e.randMu.Unlock()
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[e.rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
