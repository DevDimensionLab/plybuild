package http

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func GetJson(url string, parsed interface{}) error {
	body, err := get(url)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return err
	}

	return nil
}

func GetXml(url string, parsed interface{}) error {
	body, err := get(url)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(body, &parsed)
	if err != nil {
		return err
	}

	return nil
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(fmt.Sprintf("%s returned status code [%s]", url, resp.Status))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	defer resp.Body.Close()

	return body, nil
}

func GetJsonWithAccessToken(host string, path string, accessToken string, response interface{}) error {
	req, err := http.NewRequest("GET", host+path, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", `Bearer `+accessToken)
	req.Header.Add("Content-Type", `application/json`)

	client := &http.Client{}
	resp, err := client.Do(req)

	if nil != err {
		log.Debugln(req.Method, host+path, resp.StatusCode, err)
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Debugln(req.Method, host+path, resp.StatusCode, len(body))

	err = json.Unmarshal(body, &response)
	if nil != err {
		return err
	}

	return nil
}
