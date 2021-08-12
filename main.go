package main

import (
	"encoding/json"
	"eureka-exporter/client"
	"log"
	"net/http"
	"strings"
)

func main() {
	url := "http://172.16.0.11:9100/"
	//url := "http://139.159.217.102:9100/"
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		appID := request.URL.Query().Get("appID")

		apps, err := client.GetApps(url + "eureka/apps")
		if err != nil {
			log.Printf("get apps error: %v", err)
			return
		}
		want := filter(apps, func(app client.App) bool {
			return strings.Contains(strings.ToLower(app.Name), appID)
		})
		header := writer.Header()
		header.Add("Content-Type", "application/json; charset=utf-8")

		content, err := json.Marshal(want)
		if err != nil {
			log.Printf("marshal json error: %v", err)
			return
		}
		_, err = writer.Write(content)
		if err != nil {
			log.Printf("write content error: %v", err)
			return
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func processData(apps []client.App) []string{
	var ans []string
	for _, val := range apps {
		for _, instance := range val.Instance {
			ans= append(ans, instance.IpAddr)
		}
	}
	return ans
}

func filter(gar *client.GetAppsResp, f func(app client.App) bool) []client.App {
	var apps []client.App
	for _, val := range gar.Apps.App {
		if f(val) {
			apps = append(apps, val)
		}
	}
	return apps
}
