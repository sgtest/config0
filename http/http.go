package util

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func RespondsWithHTTPOK(url string) bool {
	resp, err := http.Get(url)
	return resp != nil && err == nil && resp.StatusCode == http.StatusOK
}

func HTTPGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get %q: http error %d", url, resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

func NewTimeoutClient(connectTimeout, readWriteTimeout, headerTimeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: TimeoutDialer(connectTimeout, readWriteTimeout),
			ResponseHeaderTimeout: headerTimeout,
		},
	}
}
