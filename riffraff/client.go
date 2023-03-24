package riffraff

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HistoryItem struct {
	Build           string `json:"build"`
	Deployer        string `json:"deployer"`
	DurationSeconds int    `json:"durationSeconds"`
	LogURL          string `json:"logURL"`
	ProjectName     string `json:"projectName"`
	Stage           string `json:"stage"`
	Status          string `json:"status"`
	Tags            struct {
		Branch            string `json:"branch"`
		Riffraff_hostname string `json:"riffraff-hostname"`
		VcsBaseURL        string `json:"vcsBaseUrl"`
		VcsCommitURL      string `json:"vcsCommitUrl"`
		VcsHeadURL        string `json:"vcsHeadUrl"`
		VcsName           string `json:"vcsName"`
		VcsRepo           string `json:"vcsRepo"`
		VcsRevision       string `json:"vcsRevision"`
		VcsTreeURL        string `json:"vcsTreeUrl"`
		VcsURL            string `json:"vcsUrl"`
	} `json:"tags"`
	Time string `json:"time"`
	UUID string `json:"uuid"`
}

type HistoryResponse struct {
	Response struct {
		CurrentPage int           `json:"currentPage"`
		Filter      []interface{} `json:"filter"`
		PageSize    int           `json:"pageSize"`
		Pages       int           `json:"pages"`
		Results     []HistoryItem `json:"results"`
		Status      string        `json:"status"`
		Total       int           `json:"total"`
	} `json:"response"`
}

type Client struct {
	url    string
	apiKey string
}

func New(url string, APIKey string) Client {
	return Client{url: url, apiKey: APIKey}
}

func (c Client) GetHistory(params HistoryParams) (HistoryResponse, error) {
	url := fmt.Sprintf("%s/api/history?key=%s&maxDaysAgo=%d&pageSize=%d&page=%d", c.url, c.apiKey, params.MaxDaysAgo, params.PageSize, params.Page)

	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		return HistoryResponse{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return HistoryResponse{}, err
	}

	historyResponse := HistoryResponse{}
	err = json.Unmarshal(data, &historyResponse)
	if err != nil {
		return HistoryResponse{}, err
	}

	return historyResponse, err
}

type HistoryParams struct {
	Page       int
	PageSize   int
	MaxDaysAgo int
}

type HistoryPaginator struct {
	HistoryParams
	client       Client
	hasMorePages bool
}

func NewHistoryPaginator(client Client, params HistoryParams) HistoryPaginator {
	return HistoryPaginator{hasMorePages: true, HistoryParams: params, client: client}
}

func (hp *HistoryPaginator) NextPage() ([]HistoryItem, error) {
	resp, err := hp.client.GetHistory(hp.HistoryParams)
	if err != nil {
		return nil, err
	}

	hp.hasMorePages = resp.Response.CurrentPage < resp.Response.Pages

	return resp.Response.Results, err
}

func (hp HistoryPaginator) HasMorePages() bool {
	return hp.hasMorePages
}
