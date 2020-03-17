# apikey
Simple API key validator middleware for [Fiber](https://github.com/gofiber/fiber).

## Install

```
go get -u github.com/fiberweb/apikey
```

## Usage

```
package main

import (
  "github.com/gofiber/fiber"
  "github.com/fiberweb/apikey"
)

func main() {
  app := fiber.New()
  
  app.Use(apikey.New(apikey.Config{Key: "secret"}))
  app.Get("/", func(c *fiber.Ctx) {
    c.Send("Ok")
  })
  app.Listen("8080")
}
```
Call the above endpint with `?key=secret` in the URL or pass the `x-api-key` header with value `secret` you will get success response.

## Default configuration
If you don't pass `apikey.Config` object to the `apikey.New()` method, it will use default configuration which will require you to specify the api key via environment variable named `API_KEY`.

```
$ API_KEY=secret ./MyGoApp
- or -
$ export API_KEY=secret
$ ./MyGoApp
```

## Config
```
type Config struct {
	// Skip this middleware
	Skip func(*fiber.Ctx) bool
  
	// Key is the api key
	Key string
  
	// ValidatorFunc is the function to validates the request
	ValidatorFunc func(*fiber.Ctx, Config) bool
}
```

#### Using Skip
`Skip` allows you to skip this middleware being executed on certain condition.
```
app.Use(apikey.New(apikey.Config{
    Skip: func(c *fiber.Ctx) bool {
      // add your logic here
      return true // returning true will skip this middleware.
    }
}))
```

#### Using ValidatorFunc
`ValidatorFunc` allows you to specify your own API key validation logic. Your custom validator function will have access to `*fiber.Ctx` and `Config`.

In this example we'll create a simple API key validator where the API key itself need to be passed as part of JSON request body called `api_key` by the client.
```
type payload struct {
    ApiKey   string `json:"api_key"`
    Command  string `json:"cmd"`
}

app.Use(apikey.New(apikey.Config{
    Key: "secret",
    ValidatorFunc: func(c *fiber.Ctx, cfg Config) bool {
        var p payload
        
        if err := json.Unmarshal([]byte(c.Body()), &p); err != nil {
            // invalid body
            return false
        }
        return cfg.Key == p.ApiKey
    }
}))
app.Post("/", func(c *fiber.Ctx) {
    c.Send("Success")
})
```
Now if you create a POST request to that endpoint with `{"api_key": "secret", "cmd": "do something"}` JSON body you'll get the success response.
