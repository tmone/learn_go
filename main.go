package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/plutov/paypal"
)

const (
	PaypalBaseUrl = "https://api-m.sandbox.paypal.com"
	ClientId      = "AdXaTMfTVHXn0G26Ay09-8I8cjJf4IBxwdcckGQA35qyo2jrbYzqhSuTRXaLCpkGYfBRbwwNRg5MmqYO"
	Secret        = "EKBtYlLAj06ONtOMETmF976Guzoi7h3h4FhCrcB_gTQ1sCSdGZKn0Rj9LbcYUIcNpbaOuFj3W0VtyEg_"
	AccessToken   = "access_token$sandbox$7ccbj4wmq4d86qbh$5f2e3aa334760ccf5208c83dc3eeec62"
)

func main() {
	c, err := paypal.NewClient(ClientId, Secret, paypal.APIBaseSandBox)
	c.SetLog(os.Stdout) // Set log to terminal stdout

	accessToken, err := c.GetAccessToken(context.Background())
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%s", accessToken.Token)

	option := paypal.TransactionSearchRequest{
		StartDate: time.Now().AddDate(0, 0, -7),
		EndDate:   time.Now(),
	}

	transactionList, err := c.ListTransactions(context.Background(), &option)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Print(transactionList)

}
