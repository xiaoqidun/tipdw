package tipdw

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 5 * time.Second}

type Body struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Result  Result `json:"result"`
}

type Result struct {
	IP       string   `json:"ip"`
	AdInfo   AdInfo   `json:"ad_info"`
	Location Location `json:"location"`
}

type AdInfo struct {
	Nation   string `json:"nation"`
	Province string `json:"province"`
	City     string `json:"city"`
	District string `json:"district"`
	Adcode   int    `json:"adcode"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// QueryIP 使用腾讯位置服务查询IP
func QueryIP(sk string, key string, ip string) (result Result, err error) {
	arg := &reqLBS{
		SK:   sk,
		Path: "/ws/location/v1/ip",
		Args: map[string]string{
			"ip":  ip,
			"key": key,
		},
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", "https://apis.map.qq.com/ws/location/v1/ip", arg.Encode()), nil)
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var bodyUnmarshal Body
	err = json.Unmarshal(body, &bodyUnmarshal)
	if err != nil {
		return
	}
	if bodyUnmarshal.Status != 0 {
		err = fmt.Errorf("resp code is %d, body is %s", bodyUnmarshal.Status, body)
		return
	}
	result = bodyUnmarshal.Result
	return
}
