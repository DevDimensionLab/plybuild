package bitbucket

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Links struct {
	Clone []struct {
		Href string `json:"href"`
		Name string `json:"name"`
	} `json:"clone"`
	Self []struct {
		Href string `json:"href"`
	} `json:"self"`
}

type Project struct {
	Key         string `json:"key"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	Type        string `json:"type"`
	Links       Links  `json:"links"`
}

type ProjectList struct {
	Size          int       `json:"size"`
	Limit         int       `json:"limit"`
	IsLastPage    bool      `json:"isLastPage"`
	Values        []Project `json:"values"`
	Start         int       `json:"start"`
	NextPageStart int       `json:"nextPageStart"`
}

type Origin struct {
	Slug          string  `json:"slug"`
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	ScmID         string  `json:"scmId"`
	State         string  `json:"state"`
	StatusMessage string  `json:"statusMessage"`
	Forkable      bool    `json:"forkable"`
	Project       Project `json:"project"`
	Public        bool    `json:"public"`
	Links         Links   `json:"links"`
}

type Repo struct {
	Slug          string  `json:"slug"`
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	ScmID         string  `json:"scmId"`
	State         string  `json:"state"`
	StatusMessage string  `json:"statusMessage"`
	Forkable      bool    `json:"forkable"`
	Project       Project `json:"project"`
	Public        bool    `json:"public"`
	Links         Links   `json:"links"`
	Origin        Origin  `json:"origin,omitempty"`
}

type ProjectRepos struct {
	Size          int    `json:"size"`
	Limit         int    `json:"limit"`
	IsLastPage    bool   `json:"isLastPage"`
	BitBucketRepo []Repo `json:"values"`
	Start         int    `json:"start"`
}

func QueryProjects(host string, personalAccessToken string) (*ProjectList, error) {
	client := &http.Client{
	}
	req, err := http.NewRequest("GET", host+"/rest/api/1.0/projects?limit=500", nil)
	req.Header.Add("Authorization", `Bearer `+personalAccessToken)
	req.Header.Add("Content-Type", `application/json`)
	resp, err := client.Do(req)
	if nil != err {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	bitBucketProjectsResponse := ProjectList{}
	err = json.Unmarshal(body, &bitBucketProjectsResponse)
	if nil != err {
		return nil, err
	}
	return &bitBucketProjectsResponse, nil
}

func Repos(host string, personalAccessToken string, projectKey string) (*ProjectRepos, error) {
	client := &http.Client{
	}

	req, err := http.NewRequest("GET", host+"/rest/api/1.0/projects/"+projectKey+"/repos?limit=1000", nil)
	req.Header.Add("Authorization", `Bearer `+personalAccessToken)
	req.Header.Add("Content-Type", `application/json`)
	resp, err := client.Do(req)
	if nil != err {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	bitBucketProjectReposResponse := ProjectRepos{}
	err = json.Unmarshal(body, &bitBucketProjectReposResponse)
	if nil != err {
		return nil, err
	}

	return &bitBucketProjectReposResponse, nil
}
