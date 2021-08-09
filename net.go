package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

func proxytest(socksstr string, url1 string, timeout time.Duration) float64 {
	var retime float64
	Proxy, err := url.Parse(socksstr)
	if err != nil {
		log.Println("[DEBUG]", err)
		return retime
	}

	method := "GET"
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(Proxy),
		},
	}

	t1 := time.Now()
	req, err := http.NewRequest(method, url1, nil)
	if err != nil {
		log.Println("[DEBUG]", err)
		return retime
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("[DEBUG]", err)
		return retime
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		retime = time.Now().Sub(t1).Seconds()
		log.Printf("[INFO] %s\t%s\t%.3f\n", socksstr, url1, retime)
	}
	return retime
}

func (i *proxyer) connectTest(time time.Duration) bool {
	_, err := net.DialTimeout("tcp", i.socksstr, time)
	if err != nil {
		log.Println("[DEBUG]", err)
		return false
	}
	return true
}
