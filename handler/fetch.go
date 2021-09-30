package handler

import (
	"cncamp_a01/constant"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	log "github.com/sirupsen/logrus"
)

func (h *handler) fetchAll() {
	var (
		wg          = &sync.WaitGroup{}
		cryptoCodes = []constant.CryptoCodeEnum{
			constant.CryptoADA, constant.CryptoBNB,
			constant.CryptoBTC, constant.CryptoETH,
		}
	)

	log.Info("updating crypto info...")
	for _, code := range cryptoCodes {
		wg.Add(1)
		go h.fetchAndUpsert(wg, code)
	}

	wg.Wait()
	log.Info("updating crypto info done")
}

type fetchResponse struct {
	Data struct {
		Symbol     string `json:"symbol"`
		Name       string `json:"name"`
		MarketData struct {
			PriceUSD float64 `json:"price_usd"`
		} `json:"market_data"`
	} `json:"data"`
}

func (h *handler) fetchAndUpsert(wg *sync.WaitGroup, code constant.CryptoCodeEnum) {
	defer wg.Done()
	const fetchAPI = "https://data.messari.io/api/v1/assets/%s/metrics"
	url := fmt.Sprintf(fetchAPI, code.String())
	resp, err := h.fetchClient.Get(url)
	if err != nil {
		log.Error(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}

	r := new(fetchResponse)
	err = json.Unmarshal(body, r)
	if err != nil {
		log.Error(err)
		return
	}

	const cryptoQuery = `
		INSERT INTO cryptos (crypto_code, price, updated_at) VALUES (?, ?, datetime('now'))
		ON CONFLICT(crypto_code) DO UPDATE SET price=?, updated_at=datetime('now');`
	_, err = h.db.Exec(cryptoQuery, code.String(), r.Data.MarketData.PriceUSD, r.Data.MarketData.PriceUSD)
	if err != nil {
		log.Error(err)
		return
	}
}
