package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ucd_url := os.Getenv("UCD_URL")
	ucd_user := os.Getenv("UCD_USER")
	ucd_password := os.Getenv("UCD_PASSWORD")
	application_id := os.Getenv("APPLICATION_ID")
	app_proccess_id := os.Getenv("APP_PROCESS_ID")
	description := os.Getenv("DESCRIPTION")
	environment_id := os.Getenv("ENVIRONMENT_ID")
	component_id := os.Getenv("COMPONENT_ID")
	version := os.Getenv("VERSION")

	if len(ucd_url) == 0 {
		log.Fatal("UCD_URL is not set")
	}
	if len(ucd_user) == 0 {
		log.Fatal("UCD_USER is not set")
	}
	if len(ucd_password) == 0 {
		log.Fatal("UCD_PASSWORD is not set")
	}
	if len(application_id) == 0 {
		log.Fatal("APPLICATION_ID is not set")
	}
	if len(app_proccess_id) == 0 {
		log.Fatal("APP_PROCESS_ID is not set")
	}
	if len(environment_id) == 0 {
		log.Fatal("ENVIRONMENT_ID is not set")
	}
	if len(component_id) == 0 {
		log.Fatal("COMPONENT_ID is not set")
	}

	log.Printf("Starting http calls for %s %s\n", component_id, version)

	bodyMap := make(map[string]interface{})
	versionArray := [1]map[string]string{
		map[string]string{
			"component": component_id,
			"version":   version}}

	bodyMap["application"] = application_id
	bodyMap["applicationProcess"] = app_proccess_id
	bodyMap["environment"] = environment_id
	if len(description) > 0 {
		bodyMap["description"] = description
	}
	bodyMap["versions"] = versionArray

	jsonIn, _ := json.Marshal(bodyMap)
	log.Println(string(jsonIn))

	client := &http.Client{}
	request, err := http.NewRequest("PUT", ucd_url+"/cli/applicationProcessRequest/request", bytes.NewReader(jsonIn))
	if err != nil {
		log.Fatal(err)
	}
	request.SetBasicAuth(ucd_user, ucd_password)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(response.Status)
	log.Println(response.StatusCode)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode >= 400 {
		log.Println(string(contents))
		log.Fatal(response.Status)
	}

	requestBody := make(map[string]string)

	log.Println(string(contents))
	err = json.Unmarshal(contents, &requestBody)

	if err != nil {
		log.Fatal(err)
	}

	var status string
	var result string

	for count := 0; status != "CLOSED" && count < 30; count++ {
		if count > 0 {
			time.Sleep(10 * time.Second)
		}
		request, err = http.NewRequest("GET", ucd_url+"/cli/applicationProcessRequest/requestStatus?request="+requestBody["requestId"], nil)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(request.URL)
		request.SetBasicAuth(ucd_user, ucd_password)
		requestResp, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(requestResp.Status)
		log.Println(requestResp.StatusCode)

		defer requestResp.Body.Close()
		contents, err = ioutil.ReadAll(requestResp.Body)
		log.Println(string(contents))

		requestStatus := make(map[string]string)

		err = json.Unmarshal(contents, &requestStatus)
		if err != nil {
			log.Fatal(err)
		}
		status = requestStatus["status"]
		result = requestStatus["result"]
	}
	if result != "SUCCEEDED" {
		log.Fatal("Request has not succeeded.  Status: " + status + ", result: " + result)
	}
	return
}
