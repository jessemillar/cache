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
	response, err := GetHttpResponse("GET", "https://statmike.michaelteamracing.com/stats/jesse", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}

func TestCacheAsStruct(t *testing.T) {
	apiResponse := testMiles{}
	_, err := GetHttpResponseAsStruct("GET", "https://statmike.michaelteamracing.com/stats/team", nil, &apiResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(apiResponse.Miles)
}
