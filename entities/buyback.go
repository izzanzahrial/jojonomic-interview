package entities

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Buyback struct {
	Norek string  `json:"norek"`
	Harga int32   `json:"harga"`
	Topup int32   `json:"topup"`
	Gram  float32 `json:"gram"`
}

func NewBuyback(norek string, harga int32, gram float32) *Buyback {
	var buyback Buyback

	buyback.populateBuyback()
	buyback.Norek = norek
	buyback.Harga = harga
	buyback.Gram = gram

	return &buyback
}

func (b *Buyback) Validate() map[string]string {
	var errors = make(map[string]string)

	if err := b.validateGramAndSaldo(); err != nil {
		errors["buyback"] = err.Error()
	}

	return errors
}

func (b *Buyback) validateGramAndSaldo() error {
	var output struct {
		Error bool `json:"error"`
		Data  struct {
			Norek int32   `json:"norek"`
			Saldo float32 `json:"saldo"`
		}
	}

	body, err := json.Marshal(map[string]string{
		"norek": b.Norek,
	})
	if err != nil {
		return err
	}

	fmt.Println("make a request to account server")
	request, err := http.NewRequest("GET", "http://localhost:8083/api/saldo", bytes.NewBuffer(body))
	if err != nil {
		return errors.New("failed to create new request")
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.New("failed to get response from saldo server")
	}
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&output); err != nil {
		return errors.New("failed to decode response")
	}
	fmt.Println("get a response from account server", output)

	if b.Gram <= output.Data.Saldo {
		return nil
	}

	return errors.New("saldo is not enough to make a buyback")
}

func (b *Buyback) populateBuyback() error {
	var output struct {
		Error bool `json:"error"`
		Data  struct {
			Buyback int32 `json:"harga_buyback"`
			Topup   int32 `json:"harga_topup"`
		}
	}

	fmt.Println("make a request to check price server")
	request, err := http.NewRequest(http.MethodGet, "http://localhost:8081/api/check-harga", nil)
	if err != nil {
		return errors.New("failed to create new request")
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.New("failed to get response from price server")
	}
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&output); err != nil {
		return errors.New("failed to decode response")
	}
	fmt.Println("get a response from check price server", output)

	if b.Harga == output.Data.Buyback {
		b.Topup = output.Data.Topup
		return nil
	}

	return errors.New("buyback is not the same as a buyback from price server")
}
