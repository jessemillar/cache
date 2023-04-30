# go-cache

I got tired of rewriting HTTP response caching logic for all my microservices so I hacked together this library that saves HTTP responses as local .txt files. `go-cache` is smart enough to return the contents of the cache file if it's only been a bit since the last call to that endpoint (cache filenames are generated as a hash of the HTTP method, URL, and headers). If the cache is expired, it remakes the call, saves the response to a file, and returns the contents either as a string or a struct. See usage examples below.

`go-cache` currently has hard-coded defaults of 60 second HTTP timeout and 60 minute cache TTL.

## Installation

```
go get github.com/jessemillar/go-cache
```

## Usage

```
import "github.com/jessemillar/go-cache"
```

### Response as String
```
response, err := GetHttpResponseAsString("GET", "https://statmike.michaelteamracing.com/stats/jesse", nil)
if err != nil {
	log.Fatal(err)
}
fmt.Println(response)
```

### Response as Struct
```
apiResponse := testMiles{}
err = GetHttpResponseAsStruct("GET", "https://statmike.michaelteamracing.com/stats/team", nil, &apiResponse)
if err != nil {
	log.Fatal(err)
}
fmt.Println(apiResponse.Miles)
```
