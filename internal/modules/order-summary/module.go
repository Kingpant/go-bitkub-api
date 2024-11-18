package ordersummary

import (
	bitkubapi "bitkub-port-avg/internal/connectors/bitkub-api"
	"fmt"
	"os"
	"strconv"
)

type IOrderSummary interface {
	GetOrderSummary(tokenSymbol string, startTimestamp *uint64) (GetOrderSummaryResponse, error)
}

type orderSummary struct {
	bitkubApiClient bitkubapi.IBitkubApiClient
}

func NewOrderSummary(bitkubApiClient bitkubapi.IBitkubApiClient) IOrderSummary {
	return &orderSummary{
		bitkubApiClient: bitkubApiClient,
	}
}

type GetOrderSummaryResponse struct {
	RateToFiatAmountBuy   map[float64]float64
	RateToTokenAmountBuy  map[float64]float64
	RateToFiatAmountSell  map[float64]float64
	RateToTokenAmountSell map[float64]float64
}

func (o *orderSummary) GetOrderSummary(tokenSymbol string, startTimestamp *uint64) (GetOrderSummaryResponse, error) {
	orderHistories, err := o.bitkubApiClient.RequestOrderHistories(tokenSymbol, startTimestamp)
	if err != nil {
		return GetOrderSummaryResponse{}, err
	}

	rateToFiatAmountBuy := map[float64]float64{}
	rateToTokenAmountBuy := map[float64]float64{}
	rateToFiatAmountSell := map[float64]float64{}
	rateToTokenAmountSell := map[float64]float64{}
	tokenRemainingAmount := 0.0
	totalInvestmentFiat := 0.0
	for _, order := range orderHistories {
		amountFloat, _ := strconv.ParseFloat(order.Amount, 64)
		feeFloat, _ := strconv.ParseFloat(order.Fee, 64)
		creditFloat, _ := strconv.ParseFloat(order.Credit, 64)
		rateFloat, _ := strconv.ParseFloat(order.Rate, 64)
		if order.Side == "buy" {
			fiatAmount := amountFloat - feeFloat + creditFloat
			tokenAmount := fiatAmount / rateFloat

			tokenRemainingAmount += tokenAmount
			totalInvestmentFiat += amountFloat

			rateToFiatAmountBuy[rateFloat] += fiatAmount
			rateToTokenAmountBuy[rateFloat] += tokenAmount
		} else {
			fiatAmount := amountFloat*rateFloat - feeFloat + creditFloat

			totalInvestmentFiat -= fiatAmount
			tokenRemainingAmount -= amountFloat

			rateToFiatAmountSell[rateFloat] += fiatAmount
			rateToTokenAmountSell[rateFloat] += amountFloat
		}
	}

	fmt.Println("====================================================")
	fmt.Printf("Token %s remaining amount: %f\n", tokenSymbol, tokenRemainingAmount)
	fmt.Printf("Total investment in fiat: %f\n", totalInvestmentFiat)
	if tokenRemainingAmount > 0 {
		fmt.Println("Average price per token: ", totalInvestmentFiat/tokenRemainingAmount)
	} else {
		fmt.Println("Average price per token: 0")
	}

	if len(rateToFiatAmountBuy) > 0 || len(rateToFiatAmountSell) > 0 {
		if err := writeToFile(tokenSymbol, rateToFiatAmountBuy, rateToTokenAmountBuy, rateToFiatAmountSell, rateToTokenAmountSell); err != nil {
			return GetOrderSummaryResponse{}, err
		}
	}

	return GetOrderSummaryResponse{
		RateToFiatAmountBuy:   rateToFiatAmountBuy,
		RateToTokenAmountBuy:  rateToTokenAmountBuy,
		RateToFiatAmountSell:  rateToFiatAmountSell,
		RateToTokenAmountSell: rateToTokenAmountSell,
	}, nil
}

func writeToFile(tokenSymbol string, rateToFiatAmountBuy, rateToTokenAmountBuy, rateToFiatAmountSell, rateToTokenAmountSell map[float64]float64) error {
	fileContent := "rate,type,fiat_amount,token_amount\n"

	for rate, fiatAmount := range rateToFiatAmountBuy {
		fileContent += fmt.Sprintf("%f,buy,%f,%f\n", rate, fiatAmount, rateToTokenAmountBuy[rate])
	}
	for rate, fiatAmount := range rateToFiatAmountSell {
		fileContent += fmt.Sprintf("%f,sell,%f,%f\n", rate, fiatAmount, rateToTokenAmountSell[rate])
	}

	// Ensure the directory exists
	if err := os.MkdirAll("./reports", os.ModePerm); err != nil {
		return err
	}

	// write to file
	file, err := os.Create(fmt.Sprintf("./reports/%s_order_summary.csv", tokenSymbol))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fileContent)
	if err != nil {
		return err
	}

	return nil
}
