# cache

I got tired of rewriting HTTP response caching logic for all my microservices so I hacked together this library that saves HTTP responses as local .txt files containing JSON data. `cache` is smart enough to return the contents of the cache file if it's only been a bit since the last call to that endpoint (cache filenames are generated as a hash of the HTTP method, URL, and headers). If the cache is expired, it remakes the call, saves the response to a file, and returns the contents either as a string or a struct.

## Installation

```
go get github.com/jessemillar/cache
```

## Usage

```
import "github.com/jessemillar/cache"
```

See `cache_test.go` for usage examples.
