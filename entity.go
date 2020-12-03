package main

type FirebaseJson struct {
	Hosting struct {
		CleanUrls bool `json:"cleanUrls"`
		Rewrites  []struct {
			Source      string `json:"source"`
			Regex       string `json:"regex"`
			Destination string `json:"destination"`
			Function    string `json:"function"`
		} `json:"rewrites"`
	} `json:"hosting"`
}
