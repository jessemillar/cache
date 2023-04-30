package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const cacheFilePrefix = "cache-"
const httpTimeout = 60 * time.Second

func main() {
	fmt.Println("Hello")
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

func cacheHttpRequestResponse(httpMethod string, url string, headers map[string]string, cacheFilename string) error {
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

	// TODO Write to file
	return json.NewDecoder(resp.Body).Decode(target)
}

func getStravaStats(userConfig *UserConfig, userDatabase string, currentData StravaAthleteStats) (StravaAthleteStats, error) {
	// Check if we recently updated Strava data
	if userConfig.LastUpdateTime.IsZero() { // We want to always fetch on launch or if we don't have a cache
		userConfig.LastUpdateTime = time.Now()
	} else {
		if time.Now().Sub(userConfig.LastUpdateTime).Minutes() < stravaCacheTTL {
			fmt.Println("Skipping getting Strava data since we recently got it")
			return currentData, nil
		} else {
			userConfig.LastUpdateTime = time.Now()
		}
	}

	err = writeStructToDatabaseFile(userDatabase, apiResponse)
	if err != nil {
		return currentData, err
	}

	return apiResponse, nil
}

func readDatabaseFile(filename string) (string, error) {
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

func readUserDatabaseFile(filename string) (StravaAthleteStats, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return StravaAthleteStats{}, err
	}

	userDatabase := StravaAthleteStats{}
	err = json.Unmarshal([]byte(data), &userDatabase)
	if err != nil {
		return StravaAthleteStats{}, err
	}

	return userDatabase, err
}
