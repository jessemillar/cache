package cache

import (
	"bytes"
	"encoding/json"
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

func cacheHttpResponse(httpMethod string, url string, headers map[string][]string) error {
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

	return writeStructToCacheFile(composeFilename(httpRequestToHash(httpMethod, url, headers)), response)
}

func getCachedHttpResponse(httpMethod string, url string, headers map[string][]string, cacheTTLOverride time.Duration) (string, error) {
	cacheFilename := composeFilename(httpRequestToHash(httpMethod, url, headers))

	if _, err := os.Stat(cacheFilename); err == nil { // Check if the cache file exists
		modifiedTime, err := getCacheFileModifiedTime(cacheFilename)
		if err != nil {
			return "", err
		}

		tempCacheTTL := cacheTTL
		if cacheTTLOverride != 0 {
			tempCacheTTL = cacheTTLOverride
		}

		if time.Since(modifiedTime) < tempCacheTTL {
			fmt.Printf("Returning cached value from %s\n", cacheFilename)
			return cacheFilename, nil
		} else {
			err = cacheHttpResponse(httpMethod, url, headers)
			if err != nil {
				return "", err
			}
		}
	} else {
		err = cacheHttpResponse(httpMethod, url, headers)
		if err != nil {
			return "", err
		}
	}

	return cacheFilename, nil
}

// HttpRequest sends an HTTP request to the specified URL and returns the HTTP response.
// The response is cached for a duration specified by cacheTTL. If cacheTTL is zero, the default cache TTL value is used.
func HttpRequest(httpMethod string, url string, headers map[string][]string, cacheTTLOverride time.Duration) (Response, error) {
	cacheFilename, err := getCachedHttpResponse(httpMethod, url, headers, cacheTTLOverride)
	if err != nil {
		return Response{}, err
	}

	return getCacheFile(cacheFilename)
}

// HttpRequestReturnStruct is the same as HttpRequest but it returns the result as a specified struct.
func HttpRequestReturnStruct(httpMethod string, url string, headers map[string][]string, cacheTTLOverride time.Duration, target interface{}) (Response, error) {
	cacheFilename, err := getCachedHttpResponse(httpMethod, url, headers, cacheTTLOverride)
	if err != nil {
		return Response{}, err
	}

	return getCacheFileAsStruct(cacheFilename, target)
}
