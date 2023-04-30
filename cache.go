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

// -- Utility functions

func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return strconv.FormatUint(uint64(h.Sum32()), 10)
}

func mapToString(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func httpRequestToString(httpMethod string, url string, headers map[string]string) string {
	return httpMethod + url + mapToString(headers)
}

func httpRequestToHash(httpMethod string, url string, headers map[string]string) string {
	return hash(httpRequestToString(httpMethod, url, headers))
}

func composeFilename(name string) string {
	return cacheFilePrefix + name + cacheFileFormat
}

// -- File IO functions

func getCacheFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func getCacheFileAsStruct(filename string, target interface{}) error {
	data, err := getCacheFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(data), &target)
	if err != nil {
		return err
	}

	return nil
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

	fmt.Printf("Writing %s with %s\n", filename, value)

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

func cacheHttpResponse(httpMethod string, url string, headers map[string]string) error {
	req, err := http.NewRequest(httpMethod, url, nil)

	// Add headers
	for key, value := range headers {
		req.Header.Add(key, value)
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

	return writeStringToCacheFile(composeFilename(httpRequestToHash(httpMethod, url, headers)), string(bytes))
}

func getCachedHttpResponse(httpMethod string, url string, headers map[string]string) (string, error) {
	cacheFilename := composeFilename(httpRequestToHash(httpMethod, url, headers))

	if _, err := os.Stat(cacheFilename); err == nil { // Check if the cache file exists
		modifiedTime, err := getCacheFileModifiedTime(cacheFilename)
		if err != nil {
			return "", err
		}

		if time.Now().Sub(modifiedTime) < cacheTTL {
			fmt.Printf("Returning cached value from %s\n", cacheFilename)
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

func GetHttpResponseAsString(httpMethod string, url string, headers map[string]string) (string, error) {
	cacheFilename, err := getCachedHttpResponse(httpMethod, url, headers)
	if err != nil {
		return "", err
	}

	return getCacheFile(cacheFilename)
}

func GetHttpResponseAsStruct(httpMethod string, url string, headers map[string]string, target interface{}) error {
	cacheFilename, err := getCachedHttpResponse(httpMethod, url, headers)
	if err != nil {
		return err
	}

	return getCacheFileAsStruct(cacheFilename, target)
}
