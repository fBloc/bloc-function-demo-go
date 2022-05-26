# bloc-function-demo-go
[中文版](/README.zh-CN.md)

This project has 2 missions：
1. Used by [run a local demo doc](https://fbloc.github.io/docs/runDemo/Go) which take you to deploy a bloc environment locally with some preset functions to let you try bloc. Those preset functions are provided by this project
2. Bloc functions under this project are used as model for developer to learn how to deploy bloc function.

## How to learn develop bloc function by this project
### Entrance
- Entrance content is in main.go:
    ```go
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

        blocClient.Run()
    }
    ```


## What you can learn from each demo function
> Please make sure you have already read [the basic doc](https://github.com/fBloc/bloc-client-go#readme) about bloc function before continue

First of all, every bloc function implemented [BlocFunctionNodeInterface](https://github.com/fBloc/bloc-client-go/blob/main/function_interface.go#L10) And have unittest code. 

Second, project's function node is developed only to demonstrate how to develope bloc function. So will not visit real data (like stock price...)

Below list each function's unique features（in common features will not be listed）：
1. sleep function. [code](/bloc_node/sleep/node.go); [unittest](/bloc_node/sleep/node_test.go)
    - this function is used to simulate a long run function。which means it should report more enough live log & progress msg to let user know the progress of it. In run() method you can see
        - report log & progress percent
        - report progress_milestone
    - see how to [define](/bloc_node/sleep/milestone.go) and report progress milestone(in run() method)
2. stock_price_monitor function. [code](/bloc_node/stock_price_monitor/node.go); [unittest](/bloc_node/stock_price_monitor/node_test.go)
    - you can see why input param's definition are nested - here we set stock's absolute_price_monitor condition into 3 input components under single one input param. e.g: [tsla_stock_code, >, 700];
    - you can see how to define a multi kind of input in IptConfig() method (string、int、 float value type & input、select frontend form type)
    - you can see how to define a multi kind of output in OptConfig() method (string、bool、json value type)
3. new_stock function. [code](/bloc_node/new_stock/node.go); [unittest](/bloc_node/new_stock/node_test.go)
    - in IptConfig() method, you can see how to build SelectOptions.
    - in OptConfig() method, new_stock_codes field defined a []string type
4. phone_sms function. [code](/bloc_node/phone_sms/node.go); [unittest](/bloc_node/phone_sms/node_test.go)
    - in IptConfig() method, you can see support multi value Select.
