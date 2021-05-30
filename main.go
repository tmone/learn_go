package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	PaypalBaseUrl = "https://api-m.sandbox.paypal.com"
	ClientId      = "AdXaTMfTVHXn0G26Ay09-8I8cjJf4IBxwdcckGQA35qyo2jrbYzqhSuTRXaLCpkGYfBRbwwNRg5MmqYO"
	Secret        = "EKBtYlLAj06ONtOMETmF976Guzoi7h3h4FhCrcB_gTQ1sCSdGZKn0Rj9LbcYUIcNpbaOuFj3W0VtyEg_"
	AccessToken   = "access_token$sandbox$7ccbj4wmq4d86qbh$5f2e3aa334760ccf5208c83dc3eeec62"
)

type Client struct {
	BaseURL    string
	apiKey     string
	HTTPClient *http.Client
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func PaypalClient() *Client {
	return &Client{
		BaseURL: PaypalBaseUrl,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

type SearchTransactionOption struct {
	TransactionId               string `url:"transaction_id"`
	TransactionType             string `url:"transaction_type"`
	TransactionStatus           string `url:"transaction_status"`
	TransactionAmount           string `url:"transaction_amount"`
	TransactionCurrency         string `url:"transaction_currency"`
	StartDate                   string `url:"start_date"`
	EndDate                     string `url:"end_date"`
	PaymentInstrumentType       string `url:"payment_instrument_type"`
	StoreId                     string `url:"store_id"`
	TerminalId                  string `url:"terminal_id"`
	Fields                      string `url:"fields"`
	BalanceAffectingRecordsOnly string `url:"balance_affecting_records_only"`
	PageSize                    int    `url:"page_size"`
	Page                        int    `url:"page"`
}

type PaypalTransaction struct {
	TransactionInfo interface{} `json:"transaction_info"`
	PayerInfo       interface{} `json:"payer_info"`
	ShippingInfo    interface{} `json:"shipping_info"`
	CartInfo        interface{} `json:"cart_info"`
	StoreInfo       interface{} `json:"store_info"`
	AuctionInfo     interface{} `json:"auction_info"`
}

type ListTransactionResponse struct {
	TransactionDetails    []PaypalTransaction `json:"transaction_details"`
	AccountNumber         string              `json:"account_number"`
	StartDate             string              `json:"start_date"`
	EndDate               string              `json:"end_date"`
	LastRefreshedDatetime string              `json:"last_refreshed_datetime"`
	Page                  int                 `json:"page"`
	TotalItems            int                 `json:"total_items"`
	TotalPages            int                 `json:"total_pages"`
	Link                  []interface{}       `json:"link"`
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AccessToken))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	fullResponse := successResponse{
		Data: v,
	}
	if err = json.NewDecoder(res.Body).Decode(&fullResponse); err != nil {
		return err
	}

	return nil
}

func (c *Client) GetListTransactions(searchTransactionOption *SearchTransactionOption) (*ListTransactionResponse, error) {
	if searchTransactionOption == nil {
		*searchTransactionOption = SearchTransactionOption{}
	}
	v, _ := query.Values(searchTransactionOption)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/reporting/transactions?%s", c.BaseURL, v.Encode()), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res := ListTransactionResponse{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func main() {
	p := PaypalClient()
	option := SearchTransactionOption{}
	transactions, err := p.GetListTransactions(&option)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%d", transactions.TotalItems)
}
