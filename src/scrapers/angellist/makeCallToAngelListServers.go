package angellist

import (
	"errors"
	"github.com/viki-org/dnscache"
	"net"
	"net/http"
	"strings"
	"time"
)

var httpTransport *http.Transport

func makeRequestToAngelListServer(method string, url string, requestQuery string, headers map[string]string, handleRedirect bool) (response *http.Response, callError error) {

	httpClient := &http.Client{
		Timeout: time.Duration(600) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !handleRedirect {
				return errors.New("found redirect")
			}

			return nil
		},
		Transport: httpTransport,
	}

	newRequestInstance, newRequestInstanceError := http.NewRequest(method, url, strings.NewReader(requestQuery))

	if newRequestInstanceError != nil {
		loggerInstance.Println(newRequestInstanceError.Error())
		callError = newRequestInstanceError
		return
	}

	newRequestInstance.Close = true
	newRequestInstance.Header.Add("DNT", "1")
	newRequestInstance.Host = "angel.co"
	newRequestInstance.Header.Add("Origin", "https://angel.co")
	newRequestInstance.Header.Add("Accept", "*/*")
	newRequestInstance.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")

	if headers != nil {
		setHeaders(newRequestInstance, headers)
	}

	response, callError = httpClient.Do(newRequestInstance)

	return
}

func setHeaders(request *http.Request, headers map[string]string) {
	for headerName, headerValue := range headers {
		request.Header.Add(headerName, headerValue)
	}
}

func init() {

	resolver := dnscache.New(time.Minute * 5)

	httpTransport = &http.Transport{
		MaxIdleConnsPerHost: 0,
		Dial: func(network string, address string) (net.Conn, error) {
			separator := strings.LastIndex(address, ":")
			ip, _ := resolver.FetchOneString(address[:separator])
			return net.Dial("tcp", ip+address[separator:])
		},
	}
}
