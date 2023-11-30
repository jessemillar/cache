package cache

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const cacheFilePrefix = "cache-"
const cacheFileFormat = ".txt"
const httpTimeout = 60 * time.Second
const cacheTTL = 60 * time.Minute

// -- Structs

type Response struct {
	StatusCode int    `json:"status"`
	Body       string `json:"body"`
}

// -- Utility functions

func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return strconv.FormatUint(uint64(h.Sum32()), 10)
}

func mapToString(m map[string][]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func httpRequestToString(httpMethod string, url string, headers map[string][]string) string {
	return httpMethod + url + mapToString(headers)
}

func httpRequestToHash(httpMethod string, url string, headers map[string][]string) string {
	return hash(httpRequestToString(httpMethod, url, headers))
}

func composeFilename(name string) string {
	return cacheFilePrefix + name + cacheFileFormat
}

// -- File IO functions

func getCacheFile(filename string) (Response, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Response{}, err
	}

	response := Response{}
	err = json.Unmarshal([]byte(data), &response)
	if err != nil {
		return Response{}, err
	}

	return response, nil
}

func getCacheFileAsStruct(filename string, target interface{}) (Response, error) {
	response, err := getCacheFile(filename)
	if err != nil {
		return Response{}, err
	}

	err = json.Unmarshal([]byte(response.Body), &target)
	if err != nil {
		return Response{}, err
	}

	return response, nil
}

func getCacheFileModifiedTime(filename string) (time.Time, error) {
	file, err := os.Stat(filename)
	if err != nil {
		return time.Time{}, err
	}

	return file.ModTime(), nil
}

func writeStringToCacheFile(filename string, value string) error {
	err := os.WriteFile(filename, []byte(value), 0666)
	if err != nil {
		return err
	}

	//fmt.Printf("Writing %s with %s\n", filename, value)
	fmt.Printf("Writing %s\n", filename)

	return nil
}

func writeStructToCacheFile(filename string, rawStruct interface{}) error {
	marshaledStruct, err := json.Marshal(&rawStruct)
	if err != nil {
		return err
	}

	return writeStringToCacheFile(filename, string(marshaledStruct))
}

// -- HTTP functions

func cacheHttpResponse(cacheFilename string, httpMethod string, url string, headers map[string][]string) error {
	req, err := http.NewRequest(httpMethod, url, nil)
	if err != nil {
		return err
	}

	// Add headers
	for key, value := range headers {
		for _, headerValue := range value {
			req.Header.Add(key, headerValue)
		}
	}

	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	response := Response{
		StatusCode: resp.StatusCode,
		Body:       string(bytes),
	}

	return writeStructToCacheFile(cacheFilename, response)
}

// checkCacheExistenceAndPermissions does what the function name says; returns (isStale, err): if err is nil, the cache file exists and we're allowed to update the cache
func checkCacheExistenceAndPermissions(cacheFilename string, cacheTTLOverride time.Duration, allowCacheUpdate bool) (bool, error) {
	permissionsErrorMessage := "permissions only allow reading from cache, not permitted to updated cache"

	if _, err := os.Stat(cacheFilename); err == nil { // Check if the cache file exists
		modifiedTime, err := getCacheFileModifiedTime(cacheFilename)
		if err != nil {
			return true, err
		}

		tempCacheTTL := cacheTTL
		if cacheTTLOverride != 0 {
			tempCacheTTL = cacheTTLOverride
		}

		if time.Since(modifiedTime) < tempCacheTTL {
			return false, nil // Cache exists and is not stale
		} else if !allowCacheUpdate {
			return false, errors.New(permissionsErrorMessage)
		} else {
			return true, nil // Cache exists in a stale state and we're allowed to update it
		}
	} else if !allowCacheUpdate {
		return false, errors.New(permissionsErrorMessage)
	}

	return true, nil // Cache doesn't exist yet and we're allowed to create it
}

func httpRequest(httpMethod string, url string, headers map[string][]string, cacheTTLOverride time.Duration, allowCacheUpdate bool) (string, error) {
	cacheFilename := composeFilename(httpRequestToHash(httpMethod, url, headers))

	isStale, err := checkCacheExistenceAndPermissions(cacheFilename, cacheTTLOverride, allowCacheUpdate)
	if err != nil {
		return "", err
	}

	if isStale {
		err = cacheHttpResponse(cacheFilename, httpMethod, url, headers)
		if err != nil {
			return "", err
		}
	}

	return cacheFilename, nil
}

// HttpRequest sends an HTTP request to the specified URL and returns the HTTP response.
// The response is cached for a duration specified by cacheTTL. If cacheTTLOverride is zero, the default cache TTL value is used.
func HttpRequest(httpMethod string, url string, headers map[string][]string, cacheTTLOverride time.Duration, allowCacheUpdate bool) (Response, error) {
	cacheFilename, err := httpRequest(httpMethod, url, headers, cacheTTLOverride, allowCacheUpdate)
	if err != nil {
		return Response{}, err
	}

	return getCacheFile(cacheFilename)
}

// HttpRequestReturnStruct is the same as HttpRequest but it returns the result as a specified struct.
// HttpRequest sends an HTTP request to the specified URL and returns the HTTP response.
// The response is cached for a duration specified by cacheTTL. If cacheTTLOverride is zero, the default cache TTL value is used.
func HttpRequestReturnStruct(httpMethod string, url string, headers map[string][]string, cacheTTLOverride time.Duration, allowCacheUpdate bool, target interface{}) (Response, error) {
	cacheFilename, err := httpRequest(httpMethod, url, headers, cacheTTLOverride, allowCacheUpdate)
	if err != nil {
		return Response{}, err
	}

	return getCacheFileAsStruct(cacheFilename, target)
}

// BasicHttpRequest makes a request with default parameters
func BasicHttpRequest(httpMethod string, url string) (Response, error) {
	return HttpRequest(httpMethod, url, nil, 0, true)
}

// BasicHttpRequestReturnStruct makes a request with default parameters
func BasicHttpRequestReturnStruct(httpMethod string, url string, target interface{}) (Response, error) {
	return HttpRequestReturnStruct(httpMethod, url, nil, 0, true, target)
}
