package main

import (
	"crypto/tls"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultTimeOut = time.Second * 5
)

type HttpClient struct {
	client   *http.Client
	username string
	password string

	addr    string
	reqeust *http.Request
}

func NewHttpClient(addr string) (*HttpClient, error) {
	//1. get http client
	client := &http.Client{
		Timeout: defaultTimeOut,
	}

	//2. check whether it is using ssl
	if !strings.HasPrefix(addr, "http") {
		addr = "http://" + addr
	}

	uaddr, err := url.Parse(addr)
	if err != nil {
		glog.Errorf("Invalid url:%v, %v", addr, err)
		return nil, err
	}
	if uaddr.Scheme == "https" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	glog.V(2).Infof("server address is: %v", addr)

	//3. generate request
	c := &HttpClient{
		client: client,
		addr:   addr,
	}
	c.generateRequest()

	return c, nil
}

// SetUser set the login user/password for the prometheus client
func (c *HttpClient) SetUser(username, password string) {
	c.username = username
	c.password = password
	c.generateRequest()
}

func (c *HttpClient) generateRequest() error {
	req, err := http.NewRequest("Get", c.addr, nil)
	if err != nil {
		glog.Errorf("Failed to generate a http.request: %v", err)
		return err
	}

	if len(c.username) > 0 {
		req.SetBasicAuth(c.username, c.password)
	}

	c.reqeust = req

	return nil
}

func (c *HttpClient) DoGet() (string, error) {
	resp, err := c.client.Do(c.reqeust)
	if err != nil {
		glog.Errorf("Failed to send http request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("Failed to read response: %v", err)
		return "", err
	}

	return string(result), nil
}

func (c *HttpClient) DoGet2(addr string) (string, error) {
	req, err := http.NewRequest("Get", addr, nil)
	if err != nil {
		glog.Errorf("Failed to generate a http.request: %v", err)
		return "", err
	}

	if len(c.username) > 0 {
		req.SetBasicAuth(c.username, c.password)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		glog.Errorf("Failed to send http request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("Failed to read response: %v", err)
		return "", err
	}

	return string(result), nil
}
