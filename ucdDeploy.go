package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log")
	"net/http"
	"os"
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

	fmt.Printf("Starting http calls for %s %s\n", component_id, version)

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
	fmt.Println(string(jsonIn))

	client := &http.Client{}
	request, err := http.NewRequest("PUT", ucd_url+"/cli/applicationProcessRequest/request", bytes.NewReader(jsonIn))
	request.SetBasicAuth(ucd_user, ucd_password)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.Status)
	fmt.Println(response.StatusCode)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(contents))
	if response.StatusCode >= 400 {
		log.Fatal(response.Status)
	}
	return
}
