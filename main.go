package main

import (
	"bloc-examples/go/stock/bloc_node/new_stock"
	"bloc-examples/go/stock/bloc_node/phone_sms"
	"bloc-examples/go/stock/bloc_node/sleep"
	"bloc-examples/go/stock/bloc_node/stock_price_monitor"

	bloc "github.com/fBloc/bloc-client-go"
)

func main() {
	clientName := "stock_go"
	blocClient := bloc.NewClient(clientName)

	blocClient.GetConfigBuilder().SetRabbitConfig(
		"blocRabbit", "blocRabbitPasswd", []string{"127.0.0.1:5672"}, "",
	).SetServer(
		"127.0.0.1", 8080,
	).BuildUp()

	stockFunctionGroup := blocClient.RegisterFunctionGroup("Stock Monitor")
	stockFunctionGroup.AddFunction("NewStockMonitor", "new stock monitor", &new_stock.NewStock{})
	stockFunctionGroup.AddFunction("PriceMonitor", "stock real time monitor", &stock_price_monitor.StockPriceMonitor{})

	NoticeFunctionGroup := blocClient.RegisterFunctionGroup("Notice")
	NoticeFunctionGroup.AddFunction("Sms", "phone short message notice", &phone_sms.SMS{})

	ToolFunctionGroup := blocClient.RegisterFunctionGroup("Tool")
	ToolFunctionGroup.AddFunction("Sleep", "do sleep between nodes", &sleep.Sleep{})

	blocClient.Run()
}
