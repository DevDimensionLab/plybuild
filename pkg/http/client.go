package http

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/internal/adapter/httpclient"
	"io"
	"net/http"
	"net/url"
	"os"
)

func GetJson(url string, parsed interface{}) error {
	return getJson(httpclient.System(), url, parsed)
}

func getJson(dependencies httpclient.Dependencies, url string, parsed interface{}) error {
	body, err := getHTTPResponse(dependencies, httpclient.Request{URL: url})
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
	return getXml(httpclient.System(), url, parsed)
}

func getXml(dependencies httpclient.Dependencies, url string, parsed interface{}) error {
	body, err := getHTTPResponse(dependencies, httpclient.Request{URL: url})
	if err != nil {
		return err
	}

	err = xml.Unmarshal(body, &parsed)
	if err != nil {
		return err
	}

	return nil
}

func getHTTPResponse(dependencies httpclient.Dependencies, request httpclient.Request) ([]byte, error) {
	resp, err := httpclient.Execute(dependencies, request)
	return responseBody(request.URL, resp, err)
}

func responseBody(url string, resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%s returned status code [%s]", url, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	return body, nil
}

func GetAuthXml(url, username, password string, parsed interface{}) error {
	return getAuthXml(httpclient.System(), url, username, password, parsed)
}

func getAuthXml(dependencies httpclient.Dependencies, url, username, password string, parsed interface{}) error {
	body, err := getHTTPResponse(dependencies, httpclient.Request{
		URL: url,
		BasicAuth: &httpclient.BasicAuth{
			Username: username,
			Password: password,
		},
	})
	if err != nil {
		return err
	}

	err = xml.Unmarshal(body, &parsed)
	if err != nil {
		return err
	}

	return nil
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
	if err != nil {
		log.Debugln(req.Method, host+path, err)
		return err
	}

	defer func() { _ = resp.Body.Close() }()
	body, _ := io.ReadAll(resp.Body)
	log.Debugln(req.Method, host+path, resp.StatusCode, len(body))

	return json.Unmarshal(body, &response)
}

func Wget(url, filepath string) error {
	return wget(wgetDependencies{
		HTTP:  httpclient.System(),
		Files: filesystem.System(),
	}, url, filepath)
}

type wgetDependencies struct {
	HTTP  httpclient.Dependencies
	Files filesystem.Dependencies
}

func wget(dependencies wgetDependencies, url, filepath string) error {
	resp, err := httpclient.Execute(dependencies.HTTP, httpclient.Request{URL: url})
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	out, err := filesystem.Create(dependencies.Files, filepath)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = filesystem.Copy(dependencies.Files, out, resp.Body)
	return err
}

func Wpost(downloadUrl, filePath string, formData url.Values) error {
	log.Debugf("downloading %s to %s with %s", downloadUrl, filePath, formData)
	resp, err := http.PostForm(downloadUrl, formData)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
