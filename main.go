package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	flag.CommandLine.Usage = func() {
		o := flag.CommandLine.Output()
		fmt.Fprintf(o, "\nUsage: %s\n", flag.CommandLine.Name())
		flag.PrintDefaults()
	}
	proxyServerPort := flag.String("port", "3000", "Frontend proxy server port")
	firebaseJsonPath := flag.String("firebase-json", "", "Path to firebase.json")
	cloudFunctionBaseUrlString := flag.String("cloud-function-base-url", "", "A Base URL of Cloud Function HTTP triggers (e.g.: http://localhost:5001/<project-id>/us-central1)")
	webAppUrlString := flag.String("web-app-url", "", "Development server (like create-react-app, parcel or webpack dev server) URL of your web app to be deployed to Firebase Hosting (e.g.: http://localhost:8080)")
	flag.Parse()
	if *firebaseJsonPath == "" || *cloudFunctionBaseUrlString == "" || *webAppUrlString == "" {
		flag.CommandLine.Usage()
		os.Exit(1)
	}

	jsonFromFile, err := ioutil.ReadFile(*firebaseJsonPath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	var firebaseJson FirebaseJson
	err = json.Unmarshal(jsonFromFile, &firebaseJson)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	webAppUrl, _ := url.Parse(*webAppUrlString)
	cloudFunctionBaseUrl, _ := url.Parse(*cloudFunctionBaseUrlString)

	fmt.Printf("Proxy server is listening on \033[1m\033[4mhttp://localhost:%s\033[0m\n", *proxyServerPort)
	http.ListenAndServe(fmt.Sprintf("localhost:%s", *proxyServerPort), ReverseProxy(webAppUrl, cloudFunctionBaseUrl, firebaseJson))
}
