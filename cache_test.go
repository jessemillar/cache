package cache

import (
	"fmt"
	"log"
	"testing"
)

type testMiles struct {
	Miles int `json:"miles"`
}

func TestCache(t *testing.T) {
	response, err := HttpRequest("GET", "https://statmike.compycore.com/mtr/stats/jesse", nil, 0, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}

func TestBasicCache(t *testing.T) {
	response, err := BasicHttpRequest("GET", "https://statmike.compycore.com/mtr/stats/jesse")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}

func TestCacheNoUpdateAllowed(t *testing.T) {
	response, err := HttpRequest("GET", "https://statmike.compycore.com/mtr/stats/poots", nil, 0, false)
	if err == nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}

func TestCacheAsStruct(t *testing.T) {
	apiResponse := testMiles{}
	_, err := HttpRequestReturnStruct("GET", "https://statmike.compycore.com/mtr/stats/team", nil, 0, true, &apiResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(apiResponse.Miles)
}

func TestBasicCacheAsStruct(t *testing.T) {
	apiResponse := testMiles{}
	_, err := BasicHttpRequestReturnStruct("GET", "https://statmike.compycore.com/mtr/stats/team", &apiResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(apiResponse.Miles)
}
