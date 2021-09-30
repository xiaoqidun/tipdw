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
	Status  int    `json:"status"`  // 状态码，0为正常，其它为异常
	Message string `json:"message"` // 对status的描述
	Result  Result `json:"result"`  // IP定位结果
}

type Result struct {
	IP       string   `json:"ip"`       // 用于定位的IP地址
	AdInfo   AdInfo   `json:"ad_info"`  // 定位行政区划信息
	Location Location `json:"location"` // 定位坐标
}

type AdInfo struct {
	Nation   string `json:"nation"`   // 国家
	Province string `json:"province"` // 省
	City     string `json:"city"`     // 市
	District string `json:"district"` // 区
	Adcode   int    `json:"adcode"`   // 行政区划代码
}

type Location struct {
	Lat float64 `json:"lat"` // 纬度
	Lng float64 `json:"lng"` // 经度
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
