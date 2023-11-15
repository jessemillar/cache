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
	response, err := HttpRequest("GET", "https://statmike.compycore.com/mtr/stats/jesse", nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}

func TestCacheAsStruct(t *testing.T) {
	apiResponse := testMiles{}
	_, err := HttpRequestReturnStruct("GET", "https://statmike.compycore.com/mtr/stats/team", nil, 0, &apiResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(apiResponse.Miles)
}
