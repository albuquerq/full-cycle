package main

import (
	"context"
	"testing"
	"time"
)

func Test_requestQuotation(t *testing.T) {
	usd, err := requestQuotation(context.Background())
	if err != nil {
		t.Error(err)
	}
	if usd.Code != "USD" && usd.Codein != "BRL" {
		t.Error("expected USD and BRL, got", usd.Code)
	}
}

func Test_setupDB(t *testing.T) {
	db, err := setupDB(context.Background(), "test.sqlite3")
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	err = db.PingContext(ctx)
	if err != nil {
		t.Error(err)
	}
}

func Test_storeQuotation(t *testing.T) {
	db, err := setupDB(context.Background(), "test.sqlite3")
	t.Cleanup(func() {
		db.Close()
	})

	if err != nil {
		t.Error(err)
	}

	u := quotation{
		Code:       "USD",
		Codein:     "BRL",
		Name:       "USDBRL",
		High:       "5.2500",
		Low:        "5.2500",
		VarBid:     "0.0000",
		PctChange:  "0.00",
		Bid:        "5.2500",
		Ask:        "5.2500",
		Timestamp:  "1709233705",
		CreateDate: "2021-03-18 15:00:00",
	}

	err = storeQuotation(context.Background(), db, u)
	if err != nil {
		t.Error(err)
	}
}
