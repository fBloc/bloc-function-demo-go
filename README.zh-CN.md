# bloc-function-demo-go
此仓库的存在有两个意义：
1. 是[部署本地demo进行试用](https://fbloc.github.io/docs/runDemo/Go)中的预置函数
2. 开发者可以通过查看此中的`bloc function`是如何开发的来进行学习开发`bloc function`

## 如何通过此仓库进行学习开发`bloc function`
### 入口
- 此项目的所有`bloc function`都可以在`main.go`中找到：
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

        // 如果你是按照[教程](https://fbloc.github.io/docs/runDemo/Go)部署的bloc环境
        // 那么下面的地址是正确的，否则请替换为你自己部署的对应地址
        blocClient.GetConfigBuilder().SetRabbitConfig(
            "blocRabbit", "blocRabbitPasswd", []string{"127.0.0.1:5672"}, "", // rabbit的地址
        ).SetServer(
            "127.0.0.1", 8080, // bloc-server的地址
        ).BuildUp()

        // 建议将你的functions按照一定的方式进行分组
        // 就像下面一样，将所有和股票监控相关的function都放到同一个function group中
        // (当然你也可以不按用处分，比如按照团队/代码仓库划/...等进行划分)
        stockFunctionGroup := blocClient.RegisterFunctionGroup("Stock Monitor")
        // 每个 AddFunction 方法都将一个 bloc function 注册到对应的group下
        // 这里你可以跳转到对应的 bloc function 对应的实现去看其到底是怎么开发的
        stockFunctionGroup.AddFunction("NewStockMonitor", "new stock monitor", &new_stock.NewStock{})
        stockFunctionGroup.AddFunction("PriceMonitor", "stock real time monitor", &stock_price_monitor.StockPriceMonitor{})

        NotifyFunctionGroup := blocClient.RegisterFunctionGroup("Notify")
        NotifyFunctionGroup.AddFunction("Sms", "phone short message notify", &phone_sms.SMS{})

        ToolFunctionGroup := blocClient.RegisterFunctionGroup("Tool")
        ToolFunctionGroup.AddFunction("Sleep", "do sleep between nodes", &sleep.Sleep{})

        blocClient.Run()
    }
    ```

### 从此项目的demo functions中，可以学到什么
> 请先确保你已经看过了关于`bloc function`的[基础文档](https://github.com/fBloc/bloc-client-go#readme)，再继续下面的内容

首先，每个 `bloc function` 都实现了 - [`BlocFunctionNodeInterface`](https://github.com/fBloc/bloc-client-go/blob/main/function_interface.go#L10) 并且编写了 单元测试

其次，此项目的 `function node` 节点开发的目的是用于演示如何进行开发`bloc function`, 故其并不会有真正的、有外部依赖的访问（比如访问股票实际数据...）

下面列出了每个 function 的特点（都相同的将不会被列出）：
1. sleep function。[代码](/bloc_node/sleep/node.go); [单元测试](/bloc_node/sleep/node_test.go)
    - 此函数是用于模拟长运行函数的。也就是说其应该尽量上报足够多的实时日志 & 进度信息，使得用户可以在 bloc 用户端就能够看到函数的运行进度。在其实现的 `run` 方法中，你可以看到：
        - 上报log & 上报进度百分比
        - 上报进度里程碑情况
2. stock_price_monitor function。[代码](/bloc_node/stock_price_monitor/node.go); [单元测试](/bloc_node/stock_price_monitor/node_test.go)
    - 你可以了解到为什么把入参设计成了两层嵌套的 - 这里我们将设置“股票价格监控”条件的入参放到了第一个参数里，且其下还有3个参数, 入参数据举例：[tsla_stock_code, >, 700]，这三个一起构成了第一个参数
    - IptConfig() 方法支持设置的多种数据类型（设置 string、int、 float 值类型 & 设置input、select 前端组件类型)
    - OptConfig() 方法支持设置的多种数据类型（设置 string、int、json 值类型）
3. new_stock function。[代码](/bloc_node/new_stock/node.go); [单元测试](/bloc_node/new_stock/node_test.go)
    - 在IptConfig() 方法中你可以看到如何构建SelectOptions
    - OptConfig() 方法中你可以看到支持多选的Select

