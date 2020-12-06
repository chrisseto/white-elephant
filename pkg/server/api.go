package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cockroachdb/cockroach-go/crdb"
	"github.com/cockroachdb/errors"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing-contrib/go-zap/log"
	"go.uber.org/zap"
)

const (
	newRoomQuery = `INSERT INTO rooms(code) VALUES (gen_random_uuid()::string) RETURNING *`
)

type Room struct {
	ID        string    `db:"id" json:"id"`
	Code      string    `db:"code" json:"code"`
	State     string    `db:"state" json:"state"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type API struct {
	db *sqlx.DB
}

func ExecuteTx(ctx context.Context, db *sqlx.DB, f func(*sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "opening tx")
	}

	return crdb.ExecuteInTx(ctx, tx, func() error {
		return f(tx)
	})
}

func (s *API) newRoom(rw http.ResponseWriter, req *http.Request) {
	var room Room
	ctx := req.Context()

	if err := ExecuteTx(ctx, s.db, func(tx *sqlx.Tx) error {
		row := tx.QueryRowxContext(ctx, newRoomQuery)
		if err := row.Err(); err != nil {
			return err
		}
		return row.StructScan(&room)
	}); err != nil {
		log.ErrorWithContext(ctx, "writing response", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(rw)
	if err := enc.Encode(&room); err != nil {
		log.ErrorWithContext(ctx, "writing response", zap.Error(err))
	}
}

func (s *API) getRoom(rw http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	var room Room
	ctx := req.Context()

	if err := ExecuteTx(ctx, s.db, func(tx *sqlx.Tx) error {
		row := tx.QueryRowxContext(ctx, `SELECT * FROM rooms WHERE id = $1`, id)
		if err := row.Err(); err != nil {
			return err
		}
		return row.StructScan(&room)
	}); err != nil {
		log.ErrorWithContext(ctx, "writing response", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	if err := enc.Encode(&room); err != nil {
		log.ErrorWithContext(ctx, "writing response", zap.Error(err))
	}
}

func (s *API) joinRoom(rw http.ResponseWriter, req *http.Request) {
	// ctx := req.Context()
}
