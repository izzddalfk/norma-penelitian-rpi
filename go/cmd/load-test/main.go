package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

const (
	addToCartFreqPerSecond   = 30
	listOfGoodsFreqPerSecond = 100
	serverAddr               = "http://192.168.1.8:9900"
)

var goodsList = map[int]float64{
	1: 3000,
	2: 1500,
	3: 1500,
	4: 2000,
	5: 1000,
	6: 2500,
	7: 2500,
}

func main() {
	duration := 5 * time.Second

	// add to cart endpoints 50r/s
	atcTargeter := newAddToCartTargeter()
	attacker := vegeta.NewAttacker()

	var wg sync.WaitGroup
	metrics := map[string]vegeta.Metrics{}

	// Add to cart and do payment
	wg.Add(1)
	go func() {
		defer wg.Done()

		rate := vegeta.Rate{
			Freq: addToCartFreqPerSecond,
			Per:  time.Second,
		}

		var atcMetrics vegeta.Metrics
		for res := range attacker.Attack(atcTargeter.newTargeter(), rate, duration, "Load test add to cart endpoint") {
			atcMetrics.Add(res)
		}
		atcMetrics.Close()

		metrics["atc"] = atcMetrics
		clearDBReq()
	}()

	// List of goods 100r/s
	wg.Add(1)
	losTargeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    fmt.Sprintf("%s/api/small/stocks", serverAddr),
	})
	go func() {
		defer wg.Done()

		rate := vegeta.Rate{
			Freq: listOfGoodsFreqPerSecond,
			Per:  time.Second,
		}

		var losMetrics vegeta.Metrics
		for res := range attacker.Attack(losTargeter, rate, duration, "Load test list of goods endpoint") {
			losMetrics.Add(res)
		}
		losMetrics.Close()

		metrics["los"] = losMetrics
	}()

	wg.Wait()

	atcMetrics := metrics["atc"]
	r := vegeta.NewTextReporter(&atcMetrics)

	fileout, err := os.Create("report.txt")
	if err != nil {
		log.Fatalf("unable to create report file due: %v", err)
	}
	defer fileout.Close()

	r.Report(fileout)

	fmt.Printf("%+v \n", metrics)
}

type addToCartTargeter struct {
	UserID *int
}

func newAddToCartTargeter() *addToCartTargeter {
	initialUserID := 1
	return &addToCartTargeter{
		UserID: &initialUserID,
	}
}

func (t *addToCartTargeter) newTargeter() vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}
		rand.Seed(time.Now().UnixNano())

		tgt.Method = http.MethodPost
		tgt.URL = fmt.Sprintf("%s/api/small/cart", serverAddr)

		randGoodsID := rand.Intn(7-1) + 1
		randTotalGoods := rand.Intn(50-1) + 1

		reqBody := addToCartReqBody{
			UserID:     *t.UserID,
			GoodsID:    randGoodsID,
			GoodsPrice: goodsList[randGoodsID],
			TotalGoods: randTotalGoods,
		}
		*t.UserID += 1

		strReqBody, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("unable to marshal request body due: %w", err)
		}

		tgt.Body = []byte(strReqBody)
		tgtReq, err := tgt.Request()
		if err != nil {
			return fmt.Errorf("unable to get target request due: %w", err)
		}

		// do request
		tgtResp, err := http.DefaultClient.Do(tgtReq)
		if err != nil {
			return fmt.Errorf("unable to make request for the target due: %w", err)
		}
		defer tgtResp.Body.Close()

		tgtRespBody, err := io.ReadAll(tgtResp.Body)
		if err != nil {
			return fmt.Errorf("unable to read response body due: %w", err)
		}
		log.Printf("[%d] Response body: %s", tgtResp.StatusCode, string(tgtRespBody))

		// do transaction when
		// 1. Successfully add goods to cart
		// 2. Cart ID is modulo 7
		if tgtResp.StatusCode == http.StatusOK {
			var atcResp respBodyAddToCart
			err = json.Unmarshal(tgtRespBody, &atcResp)
			if err != nil {
				return fmt.Errorf("unable to unmarshal response due: %w", err)
			}
			if atcResp.Data.CartID%7 == 0 {
				tgt.URL = fmt.Sprintf("%s/api/small/pay", serverAddr)

				payReqBody := payReqBody{
					CartID:      atcResp.Data.CartID,
					TotalAmount: atcResp.Data.TotalAmount,
				}

				strPayReqBody, err := json.Marshal(payReqBody)
				if err != nil {
					return fmt.Errorf("unable to marshal request body for transactions due: %w", err)
				}

				tgt.Body = []byte(strPayReqBody)
			}
		}

		return nil
	}
}

func clearDBReq() {
	resp, err := http.Post(fmt.Sprintf("%s/clear-db", serverAddr), "application/json", nil)
	if err != nil {
		log.Printf("failed clear db: %v", err)
	} else {
		log.Printf("Response clear db: %+v", resp.StatusCode)
	}
}
