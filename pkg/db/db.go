package db

import (
	"database/sql"

	"github.com/xmtp/xmtpd/pkg/db/queries"
)

type handlerConfig struct {
	readReplica *sql.DB
}

type HandlerOption func(*handlerConfig)

func WithReadReplica(db *sql.DB) HandlerOption {
	return func(cfg *handlerConfig) {
		cfg.readReplica = db
	}
}

// Handler eases working with two databases - a read-write and read-only database. It mitigates the possibility of a component
// attempting a write to DB, not knowing it received a handle to a read-only SQL DB. Handler also makes the query intent explicit.
// The handler will correctly route the request to the appropriate DB. It also eases the transition if some part of the code used
// to do read-only access and later needs to write data.
type Handler struct {
	// Handle to read-write DB.
	write *sql.DB

	// Handle to read-only DB. Preferred for reads, if available.
	read *sql.DB

	// NOTE: This is potentially just overhead since it's trivial to create queries in other places.
	query     *queries.Queries
	readQuery *queries.Queries
}

// NewDBHandler creates a new database handler with two database connections - a read-write and a read one.
// If there's no exclusive read replica it can be ommitted and the write replica will be used.
func NewDBHandler(db *sql.DB, options ...HandlerOption) *Handler {
	var cfg handlerConfig
	for _, opt := range options {
		opt(&cfg)
	}

	handler := &Handler{
		write: db,
		query: queries.New(db),
	}

	if cfg.readReplica != nil {
		handler.read = cfg.readReplica
		handler.readQuery = queries.New(cfg.readReplica)
	}

	return handler
}

func (h *Handler) DB() *sql.DB {
	return h.Write()
}

func (h *Handler) Write() *sql.DB {
	return h.write
}

func (h *Handler) Read() *sql.DB {
	if h.read != nil {
		return h.read
	}

	return h.write
}

func (h *Handler) Query() *queries.Queries {
	return h.WriteQuery()
}

func (h *Handler) WriteQuery() *queries.Queries {
	return h.query
}

func (h *Handler) ReadQuery() *queries.Queries {
	if h.readQuery != nil {
		return h.readQuery
	}

	return h.query
}
