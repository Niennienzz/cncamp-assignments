package handler

import (
	"cncamp_a01/httpserver/constant"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	upsert := true
	res := h.mongoDB.Collection(cryptosCol).FindOneAndUpdate(
		context.Background(),
		bson.M{"crypto_code": code},
		bson.M{
			"$set": bson.M{
				"price":      r.Data.MarketData.PriceUSD,
				"updated_at": time.Now(),
			},
		},
		&options.FindOneAndUpdateOptions{Upsert: &upsert},
	)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Infof("created crypto %s", code)
			return
		}
		log.Error(err)
		return
	}
	log.Infof("updated crypto %s", code)
}
