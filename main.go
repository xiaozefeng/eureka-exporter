package main

import (
	"encoding/json"
	"eureka-exporter/client"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	url := os.Getenv("EUREKA_URL")
	log.Println("url:", url)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		appID := request.URL.Query().Get("appID")

		apps, err := client.GetApps(url + "eureka/apps")
		if err != nil {
			log.Printf("get apps error: %v", err)
			return
		}

		filtered := filter(apps, func(app client.App) bool {
			return strings.Contains(strings.ToLower(app.Name), appID)
		})
		wrapped := wrap(filtered)
		header := writer.Header()
		header.Add("Content-Type", "application/json; charset=utf-8")

		content, err := json.Marshal(wrapped)
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
	log.Println("server at: 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func wrap(apps []client.App) []AppWrapper {
	var ans []AppWrapper
	for _, val := range apps {
		var ap AppWrapper
		ap.Name = val.Name
		for _, instance := range val.Instance {
			ap.Ips = append(ap.Ips, instance.IpAddr)
		}
		ans = append(ans, ap)
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

type AppWrapper struct {
	Name string   `json:"name,omitempty"`
	Ips  []string `json:"ips,omitempty"`
}
