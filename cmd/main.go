package main

import (
	"strings"

	"bitkub-port-avg/internal/config"
	bitkubapi "bitkub-port-avg/internal/connectors/bitkub-api"
	ordersummary "bitkub-port-avg/internal/modules/order-summary"
)

func main() {
	config := config.NewConfig()
	bitkubapiClient := bitkubapi.NewBitkubApiClient(config.BitkubApiBaseUrl, config.BitkubApiKey, config.BitkubApiSecret)

	orderSummaryModule := ordersummary.NewOrderSummary(bitkubapiClient)

	tokenSymbols := strings.Split(config.Tokens, ",")
	for _, tokenSymbol := range tokenSymbols {
		_, err := orderSummaryModule.GetOrderSummary(tokenSymbol, config.StartTimestamp)
		if err != nil {
			panic(err)
		}
	}
}
