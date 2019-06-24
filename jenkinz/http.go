package jenkinz

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

var Jenkinz = &http.Client{
	Timeout: time.Second * 10,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

func Get(url string, credentials string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if len(credentials) > 0 {
		auth := base64.StdEncoding.EncodeToString([]byte(credentials))
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
	}
	resp, err := Jenkinz.Do(req)
	if err != nil {
		return nil, err
	}
	// let the calling function close the body.
	return resp, nil
}
