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
	response, err := HttpRequest("GET", "https://raw.githubusercontent.com/jessemillar/static-json/main/cache-test.json", nil, 0, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}

func TestBasicHttpCache(t *testing.T) {
	response, err := BasicHttpRequest("GET", "https://raw.githubusercontent.com/jessemillar/static-json/main/cache-test.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}

func TestHttpCacheNoUpdateAllowed(t *testing.T) {
	response, err := HttpRequest("GET", "https://raw.githubusercontent.com/jessemillar/static-json/main/cache-test.json", nil, 0, false)
	if err == nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}

func TestHttpCacheAsStruct(t *testing.T) {
	apiResponse := testMiles{}
	_, err := HttpRequestReturnStruct("GET", "https://raw.githubusercontent.com/jessemillar/static-json/main/cache-test.json", nil, 0, true, &apiResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(apiResponse.Miles)
}

func TestBasicHttpCacheAsStruct(t *testing.T) {
	apiResponse := testMiles{}
	_, err := BasicHttpRequestReturnStruct("GET", "https://raw.githubusercontent.com/jessemillar/static-json/main/cache-test.json", &apiResponse)
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
