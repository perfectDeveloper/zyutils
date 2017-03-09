package httpman

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	httpUrlPrefix = "http://"
)

var (
	httpClient *http.Client
)

func init() {
	httpClient = &http.Client{
		Transport: &http.Transport{
			Dial: localDial,
			ResponseHeaderTimeout: 60 * time.Second,
			DisableKeepAlives:     false,
			MaxIdleConnsPerHost:   10000,
		},
	}
}

/*
func main() {
	huststoreBaseUrl := "s17.basic.shbt.qihoo.net:8082/"
	data := []byte{1, 2, 3}
	code := rawPost(huststoreBaseUrl+"hset?tb=testHttp&key=fxxxxxxx", nil, "huststore", "huststore", data)
	//code, _ := rawGet(huststoreBaseUrl+"hget?tb=testHttp&key=123", nil, "huststore", "huststore")
	if code != http.StatusOK {
		fmt.Printf("http.Get error : %#v\n", code)
		return
	}

	code, body := rawGet(huststoreBaseUrl+"hget?tb=testHttp&key=fxxxxxxx", nil, "huststore", "huststore")
	if code != http.StatusOK {
		fmt.Printf("http.Get error : %#v", code)
		return
	}
	fmt.Printf("resp : %#v\n", string(body))
}
*/
func RawGet(urlStr string, header http.Header, user string, pwd string) (int, []byte) {
	if !strings.HasPrefix(urlStr, httpUrlPrefix) {
		var buff bytes.Buffer
		buff.WriteString(httpUrlPrefix)
		buff.WriteString(urlStr)
		urlStr = buff.String()
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return 0, nil
	}

	if header != nil {
		req.Header = header
	}
	req.SetBasicAuth(user, pwd)

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil
	}

	return resp.StatusCode, body
}

func RawPost(urlStr string, header http.Header, user string, pwd string, data []byte) int {
	if !strings.HasPrefix(urlStr, httpUrlPrefix) {
		var buff bytes.Buffer
		buff.WriteString(httpUrlPrefix)
		buff.WriteString(urlStr)
		urlStr = buff.String()
	}

	body := ioutil.NopCloser(bytes.NewReader(data))
	req, err := http.NewRequest("POST", urlStr, body)
	if err != nil {
		return 0
	}

	if header != nil {
		req.Header = header
	}
	req.SetBasicAuth(user, pwd)

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func localDial(network, addr string) (net.Conn, error) {
	dial := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	return dial.Dial(network, addr)
}
