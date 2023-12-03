package cache

import (
	"testing"
)

type testResponse struct {
	Test struct {
		Value int `json:"value"`
	} `json:"test"`
}

func TestHttpCache(t *testing.T) {
	response, err := HttpRequest("GET", "https://raw.githubusercontent.com/jessemillar/static-json/main/cache-test.json", nil, 0, true)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Error("Response status code is incorrect")
	}

	if response.Body != "{\n    \"test\": {\n        \"value\": 1\n    }\n}\n" {
		t.Error("Response body is incorrect")
	}
}

func TestBasicHttpCache(t *testing.T) {
	response, err := BasicHttpRequest("GET", "https://raw.githubusercontent.com/jessemillar/static-json/main/cache-test.json")
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Error("Response status code is incorrect")
	}

	if response.Body != "{\n    \"test\": {\n        \"value\": 1\n    }\n}\n" {
		t.Error("Response body is incorrect")
	}
}

// TestHttpCacheNoUpdateAllowed tries to get a file that we haven't retrieved before and it should fail because we don't allow it to update the cache by saving a file to disk
func TestHttpCacheNoUpdateAllowed(t *testing.T) {
	response, err := HttpRequest("GET", "https://raw.githubusercontent.com/jessemillar/static-json/main/cache-test-nonexistent.json", nil, 0, false)
	if err == nil {
		t.Error(err)
	}

	if response.StatusCode != 0 {
		t.Error("Response status code is incorrect")
	}

	if response.Body != "" {
		t.Error("Response body is incorrect")
	}
}

func TestHttpCacheAsStruct(t *testing.T) {
	apiResponse := testResponse{}
	err := HttpRequestReturnStruct("GET", "https://raw.githubusercontent.com/jessemillar/static-json/main/cache-test.json", nil, 0, true, &apiResponse)
	if err != nil {
		t.Error(err)
	}

	if apiResponse.Test.Value != 1 {
		t.Error("Didn't get the expected value")
	}
}

func TestBasicHttpCacheAsStruct(t *testing.T) {
	apiResponse := testResponse{}
	err := BasicHttpRequestReturnStruct("GET", "https://raw.githubusercontent.com/jessemillar/static-json/main/cache-test.json", &apiResponse)
	if err != nil {
		t.Error(err)
	}

	if apiResponse.Test.Value != 1 {
		t.Error("Didn't get the expected value")
	}
}

func TestCache(t *testing.T) {
	cacheValue := Response{}
	_, err := GetCacheAndStaleness("cache-test.txt", 0, true, &cacheValue)
	if err != nil {
		t.Error(err)
	}

	if cacheValue.StatusCode != 200 {
		t.Error("Response status code is incorrect")
	}

	if cacheValue.Body != "Test value" {
		t.Error("Response body is incorrect")
	}
}

func TestCacheAsStruct(t *testing.T) {
	cacheValue := testResponse{}
	_, err := GetCacheAndStalenessReturnStruct("cache-test-struct.txt", 0, true, &cacheValue)
	if err != nil {
		t.Error(err)
	}

	if cacheValue.Test.Value != 1 {
		t.Error("Didn't get the expected value")
	}
}

func TestGetCacheFileAsStruct(t *testing.T) {
	cacheValue := testResponse{}
	err := getCacheFileAsStruct("cache-test-struct.txt", &cacheValue)
	if err != nil {
		t.Error(err)
	}

	if cacheValue.Test.Value != 1 {
		t.Error("Didn't get the expected value")
	}
}
