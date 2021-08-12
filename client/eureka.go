package client

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func GetApps(url string) (*GetAppsResp, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	c := http.DefaultClient
	response, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	body := response.Body
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(body)

	content, err := ioutil.ReadAll(body)
	log.Printf("%s", content)
	if err != nil {
		return nil, err
	}
	var ans GetAppsResp
	err = json.Unmarshal(content, &ans)
	if err != nil {
		return nil, err
	}
	return &ans, nil
}

type GetAppsResp struct {
	Apps Apps `json:"applications"`
}

type Apps struct {
	App []App `json:"application"`
}

type App struct {
	Name     string     `json:"name"`
	Instance []Instance `json:"instance"`
}

type Instance struct {
	InstanceId string `json:"instanceId"`
	HostName   string `json:"hostName"`
	Host       string `json:"host"`
	IpAddr     string `json:"ipAddr"`
	Status     string `json:"status"`
}
