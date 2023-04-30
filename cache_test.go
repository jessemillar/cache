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
	body, response, err := GetHttpResponseAsString("GET", "https://statmike.michaelteamracing.com/stats/jesse", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(body, response)

	apiResponse := testMiles{}
	response, err = GetHttpResponseAsStruct("GET", "https://statmike.michaelteamracing.com/stats/team", nil, &apiResponse)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(apiResponse.Miles)
}
