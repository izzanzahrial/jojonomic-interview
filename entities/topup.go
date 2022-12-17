package entities

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Topup struct {
	Norek   string  `json:"norek"`
	Harga   int32   `json:"harga"`
	Buyback int32   `json:"buyback"`
	Gram    float32 `json:"gram"`
}

func NewTopup(norek string, harga int32, gram float32) *Topup {
	return &Topup{
		Norek: norek,
		Harga: harga,
		Gram:  gram,
	}
}

func (t *Topup) Validate() map[string]string {
	var errors = make(map[string]string)

	if err := t.validateHargaAndAddBuyback(); err != nil {
		errors["harga"] = err.Error()
	}

	return errors
}

func (t *Topup) validateHargaAndAddBuyback() error {
	var output struct {
		Data struct {
			Buyback int32 `json:"harga_buyback"`
			Topup   int32 `json:"harga_topup"`
		}
		Error bool `json:"error"`
	}

	fmt.Println("make a request to check price server")
	request, err := http.NewRequest("GET", "http://localhost:8081/api/check-harga", nil)
	if err != nil {
		return fmt.Errorf("failed to create new request %s", err.Error())
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

	if t.Harga == output.Data.Topup {
		t.Buyback = output.Data.Buyback
		return nil
	}

	return errors.New("price is not the same as a topup price")
}
