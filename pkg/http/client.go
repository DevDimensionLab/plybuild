package http

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
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
