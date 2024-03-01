package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config{
		listenAddr:            ":8080",
		database:              "database.sqlite3",
		getQuotationTimeout:   200 * time.Millisecond,
		storeQuotationTimeout: 10 * time.Millisecond,
	}

	db, err := setupDB(context.Background(), cfg.database)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	s := newServer(cfg, db)

	err = http.ListenAndServe(cfg.listenAddr, s)
	if err != nil {
		panic(err)
	}
}

type config struct {
	listenAddr            string
	database              string
	getQuotationTimeout   time.Duration
	storeQuotationTimeout time.Duration
}

type server struct {
	cfg config
	log *slog.Logger
	db  *sql.DB
	*http.ServeMux
}

func newServer(cfg config, db *sql.DB) *server {
	s := &server{
		cfg:      cfg,
		db:       db,
		log:      slog.Default(),
		ServeMux: http.NewServeMux(),
	}

	s.HandleFunc("/cotacao", s.getQuotationHandler)

	return s
}

func (s *server) getQuotationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	select {
	case <-ctx.Done():
		s.log.Error("client-side timeout")
	default:
		gqCtx, gqCancel := context.WithTimeout(ctx, s.cfg.getQuotationTimeout)
		defer gqCancel()

		q, err := requestQuotation(gqCtx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				s.log.Error("timeout requesting quotation", "error", err)
				w.WriteHeader(http.StatusGatewayTimeout)
				return
			}
			s.log.Error("error requesting quotation", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		sqCtx, sqCancel := context.WithTimeout(ctx, s.cfg.storeQuotationTimeout)
		defer sqCancel()

		err = storeQuotation(sqCtx, s.db, q)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				s.log.Error("timeout storing quotation", "error", err)
				w.WriteHeader(http.StatusGatewayTimeout)
				return
			}
			s.log.Error("error storing quotation", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		type response struct {
			BID string `json:"bid"`
		}

		data, err := json.Marshal(response{
			BID: q.Bid,
		})
		if err != nil {
			s.log.Error("error encoding response", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

type quotation struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func requestQuotation(ctx context.Context) (quotation, error) {
	const targetURL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, http.NoBody)
	if err != nil {
		return quotation{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return quotation{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return quotation{}, errors.New("error getting quotation")
	}

	var data struct {
		USDBRL quotation `json:"USDBRL"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return quotation{}, err
	}

	return data.USDBRL, nil
}

func setupDB(ctx context.Context, name string) (*sql.DB, error) {
	const migration = `
	CREATE TABLE IF NOT EXISTS usd_brl (
		code       	TEXT NOT NULL,
		codein     	TEXT NOT NULL,
		name       	TEXT NOT NULL,
		high       	TEXT NOT NULL,
		low        	TEXT NOT NULL,
		varBid     	TEXT NOT NULL,
		pct_change 	TEXT NOT NULL,
		bid        	TEXT NOT NULL,
		ask        	TEXT NOT NULL,
		timestamp  	TEXT NOT NULL,
		create_date TEXT NOT NULL
	);
	`

	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx, migration)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func storeQuotation(ctx context.Context, db *sql.DB, q quotation) error {
	const cmd = `
	INSERT INTO usd_brl (
		code,
		codein,
		name,
		high,
		low,
		varBid,
		pct_change,
		bid,
		ask,
		timestamp,
		create_date
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := db.PrepareContext(ctx, cmd)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		q.Code,
		q.Codein,
		q.Name,
		q.High,
		q.Low,
		q.VarBid,
		q.PctChange,
		q.Bid,
		q.Ask,
		q.Timestamp,
		q.CreateDate,
	)
	if err != nil {
		return err
	}
	return nil
}
