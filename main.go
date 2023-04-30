package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
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

func main() {
	fmt.Println("Hello")
	fmt.Println(hash("Hello"))
	fmt.Println(hash("Hello."))

	err := cacheHttpResponse("GET", "https://statmike.michaelteamracing.com/stats/jesse", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return strconv.FormatUint(uint64(h.Sum32()), 10)
}

func composeFilename(name string) string {
	return cacheFilePrefix + name + cacheFileFormat
}

func getCacheFileAsStruct(filename string, target interface{}) (interface{}, error) {
	cacheFileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(cacheFileContents), &target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

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

	return writeStringToDatabaseFile(composeFilename(hash(url)), string(bytes))
}

func getStravaStats() error {
	//if time.Now().Sub(time.Now()).Minutes() < cacheTTL { // TODO Compare now() to the file modification time
	if "poots" == "poots" {
		fmt.Println("Skipping getting Strava data since we recently got it")
		return nil
	} else {
		// TODO Do the HTTP request
	}

	return nil
}

func readCacheFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Printf("Read %s from %s\n", string(data), filename)

	return string(data), nil
}

func getCacheFileModifiedTime(filename string) (time.Time, error) {
	file, err := os.Stat(filename)
	if err != nil {
		return time.Time{}, err
	}

	return file.ModTime(), nil
}

func writeStringToDatabaseFile(filename string, value string) error {
	err := os.WriteFile(filename, []byte(value), 0666)
	if err != nil {
		return err
	}

	fmt.Printf("Updating %s to %s\n", filename, value)

	return nil
}

func writeStructToDatabaseFile(filename string, rawStruct interface{}) error {
	marshaledStruct, err := json.Marshal(&rawStruct)
	if err != nil {
		return err
	}

	return writeStringToDatabaseFile(filename, string(marshaledStruct))
}
