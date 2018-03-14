# goconsulprops
Consul key/value versioned string properties that can be refreshed.

I started using consul with spring to pull key/values into a configuration properties class in java. I wanted to keep the functionality in Go. You can use the consul UI to set key/values, then your Go program can periodically Refresh() during runtime - this allows you to manually adjust your apps/services while running (change configurations without restarting).

I found this useful for things like:
* adding/removing consumers
* changing logging levels
* graceful shutdowns
* ... 

Simple example:
```go
consulAddress := "localhost:8500"
consulPrefix := "config/GoTest/app"

appProps := goconsulprops.NewProperties(consulAddress, consulPrefix)

fmt.Printf("app consumerCount = %v, (version:%v)\n", appProps.GetValue("consumerCount"), appProps.GetVersion("consumerCount"))
```

Output:
```
2018/03/14 16:44:36 [goconsulprops] set app.consumerCount: 22 (version: 13041509)
2018/03/14 16:44:36 [goconsulprops] set app.logFile: application.log (version: 13032513)
2018/03/14 16:44:36 [goconsulprops] set app.logLevel: INFO (version: 13032525)
details for consumerCount = value:22, version:13041509
```
