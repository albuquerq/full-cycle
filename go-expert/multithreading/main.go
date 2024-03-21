package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

var cep = flag.String("cep", "01153000", "CEP to search")

func main() {
	flag.Parse()

	select {
	case msg := <-viacep(*cep):
		slog.Info("found", "api", "viacep", "address", msg)
	case msg := <-brasilapi(*cep):
		slog.Info("found", "api", "brasilapi", "address", msg)
	case <-time.After(time.Second):
		slog.Error("timeout")
	}
}

type msg struct {
	Cep          string
	State        string
	City         string
	Neighborhood string
	Street       string
}

func (m msg) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("cep", m.Cep),
		slog.String("state", m.State),
		slog.String("city", m.City),
		slog.String("neighborhood", m.Neighborhood),
		slog.String("street", m.Street),
	)
}

func get[MSG any](url string) (msg MSG, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return msg, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return msg, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&msg)

	return msg, err
}

func viacep(cep string) <-chan msg {
	type viacepResponse struct {
		Cep         string `json:"cep"`
		Logradouro  string `json:"logradouro"`
		Complemento string `json:"complemento"`
		Bairro      string `json:"bairro"`
		Localidade  string `json:"localidade"`
		Uf          string `json:"uf"`
		Ibge        string `json:"ibge"`
		Gia         string `json:"gia"`
		Ddd         string `json:"ddd"`
		Siafi       string `json:"siafi"`
	}

	ch := make(chan msg)

	go func() {
		b, err := get[viacepResponse]("https://viacep.com.br/ws/" + cep + "/json/")
		if err != nil {
			slog.Default().Error("getting from viacep", "err", err)
			return
		}

		ch <- msg{
			Cep:          b.Cep,
			State:        b.Uf,
			City:         b.Localidade,
			Neighborhood: b.Bairro,
			Street:       b.Logradouro,
		}
	}()

	return ch
}

func brasilapi(cep string) <-chan msg {
	type brasilResponse struct {
		Cep          string `json:"cep"`
		State        string `json:"state"`
		City         string `json:"city"`
		Neighborhood string `json:"neighborhood"`
		Street       string `json:"street"`
		Service      string `json:"service"`
	}

	ch := make(chan msg)

	go func() {
		b, err := get[brasilResponse]("https://brasilapi.com.br/api/cep/v1/" + cep)
		if err != nil {
			slog.Default().Error("getting from brasilapi", "err", err)
			return
		}

		ch <- msg{
			Cep:          b.Cep,
			State:        b.State,
			City:         b.City,
			Neighborhood: b.Neighborhood,
			Street:       b.Street,
		}
	}()

	return ch
}
