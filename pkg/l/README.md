# Logger - [zap](https://github.com/uber-go/zap/).Logger wrapper

## Enable debug level logger

```
export TOAA_LOG_DEBUG="*"					// enable debug level for all package
export TOAA_LOG_DEBUG="payment-gateway-service/*"	// enable debug level for all package start with identity/bank
export TOAA_LOG_DEBUG="*/service"
export TOAA_LOG_DEBUG="payment-gateway-service/cmd"

```

## Enable change log level at runtime

ServeHTTP() used to handle HTTP request to change log level at runtime

```
m.HandleFunc("/log/level", l.ServeHTTP)

// request and response message
// Name is directory name of file call New()
type payload struct {
		Name  string     `json:"name"`
		Level *zap.Level `json:"level,omitempty"`
}
```

Support 2 methods: GET and PUT.

GET: get log level of all loggers

PUT: change log level of a package name
