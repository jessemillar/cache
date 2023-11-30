package cache

import (
	"fmt"
	"log"
	"testing"
)

type testMiles struct {
	Miles int `json:"miles"`
}

func TestHttpCache(t *testing.T) {
	response, err := HttpRequest("GET", "https://statmike.compycore.com/mtr/stats/jesse", nil, 0, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}

func TestBasicHttpCache(t *testing.T) {
	response, err := BasicHttpRequest("GET", "https://statmike.compycore.com/mtr/stats/jesse")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}

func TestHttpCacheNoUpdateAllowed(t *testing.T) {
	response, err := HttpRequest("GET", "https://statmike.compycore.com/mtr/stats/poots", nil, 0, false)
	if err == nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}

func TestHttpCacheAsStruct(t *testing.T) {
	apiResponse := testMiles{}
	_, err := HttpRequestReturnStruct("GET", "https://statmike.compycore.com/mtr/stats/team", nil, 0, true, &apiResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(apiResponse.Miles)
}

func TestBasicHttpCacheAsStruct(t *testing.T) {
	apiResponse := testMiles{}
	_, err := BasicHttpRequestReturnStruct("GET", "https://statmike.compycore.com/mtr/stats/team", &apiResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(apiResponse.Miles)
}

func TestCache(t *testing.T) {
	cacheValue, isStale, err := GetCacheAndStaleness("cache-test.txt", 0, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cacheValue, isStale)
}

func TestCacheAsStruct(t *testing.T) {
	type testStruct struct {
		Test string `json:"test"`
	}

	cacheValue := testStruct{}
	isStale, err := GetCacheAndStalenessReturnStruct("cache-test-struct.txt", 0, true, &cacheValue)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cacheValue, isStale)
}
