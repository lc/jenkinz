package jenkinz

import (
	"crypto/tls"
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
