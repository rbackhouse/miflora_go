package http

import (
	"bytes"
	"fmt"
	"net/http"

	logger "github.com/sirupsen/logrus"

	"potpie.org/miflora/src/config"
)

type httpData struct {
	config config.HttpConfig
}

func sendHttpRequest(config config.HttpConfig, payloadBuf *bytes.Buffer, urlSuffix string) {
	url := fmt.Sprintf("%s:%d/%s", config.HttpHost, config.HttpPort, urlSuffix)
	logger.Infof("Sending HTTP Request to : %s", url)
	req, _ := http.NewRequest("POST", url, payloadBuf)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, e := client.Do(req)
	if e != nil {
		logger.Errorf("Failed to send http request to : %s err: %s", url, e)
	}

	defer res.Body.Close()
}
