package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	timeout  = 300 * time.Millisecond
	target   = "http://localhost:8080/cotacao"
	fileName = "cotacao.txt"
)

func main() {
	logger := slog.Default()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, http.NoBody)
	if err != nil {
		logger.Error("error creating http request", "error", err)
		return
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Error("timeout requesting quotation")
			return
		}
		logger.Error("error requesting quotation", "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("error requesting quotation", "status", resp.StatusCode)
		return
	}

	var data struct {
		BID string `json:"bid"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		logger.Error("error decoding response", "error", err)
		return
	}

	logger.Info("Cotação recebida do servidor", "bid", data.BID)

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logger.Error("error opening file", "error", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "Dólar: %s\n", data.BID)

	logger.Info("Cotação armazenada com sucesso", "file", fileName)
}
