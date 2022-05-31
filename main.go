package main

import (
	"stock/bloc_node/new_stock"
	"stock/bloc_node/phone_sms"
	"stock/bloc_node/sleep"
	"stock/bloc_node/stock_price_monitor"

	bloc "github.com/fBloc/bloc-client-go"
)

func main() {
	clientName := "stock_go"
	blocClient := bloc.NewClient(clientName)

	// below address is address if you deploy bloc by tutorial https://fbloc.github.io/docs/runDemo/Go
	// if you deploy bloc by yourself, you should change it to your own address
	blocClient.GetConfigBuilder().SetRabbitConfig(
		"blocRabbit", "blocRabbitPasswd", []string{"127.0.0.1:5672"}, "", // bloc use rabbitMQ address
	).SetServer(
		"127.0.0.1", 8080, // bloc-server address
	).BuildUp()

	// Recommend group bloc functions, like below put all stock monitor about bloc functions into one group
	stockFunctionGroup := blocClient.RegisterFunctionGroup("Stock Monitor")
	// Every AddFunction method added a bloc function to register
	// so you can jump to each of it to learn how the bloc function is developed
	stockFunctionGroup.AddFunction("NewStockMonitor", "new stock monitor", &new_stock.NewStock{})
	stockFunctionGroup.AddFunction("PriceMonitor", "stock real time monitor", &stock_price_monitor.StockPriceMonitor{})

	NotifyFunctionGroup := blocClient.RegisterFunctionGroup("Notify")
	NotifyFunctionGroup.AddFunction("Sms", "phone short message notify", &phone_sms.SMS{})

	ToolFunctionGroup := blocClient.RegisterFunctionGroup("Tool")
	ToolFunctionGroup.AddFunction("Sleep", "do sleep between nodes", &sleep.Sleep{})

	// you can start multi bloc client to run concurrently
	for i := 0; i < 1; i++ {
		go blocClient.Run()
	}
	forever := make(chan struct{})
	<-forever
}
