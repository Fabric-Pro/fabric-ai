package jina

// see https://jina.ai for more information

import (
	"fmt"
	"io"
	"net/http"

	"github.com/danielmiessler/fabric/internal/plugins"
)

type Client struct {
	*plugins.PluginBase
	ApiKey *plugins.SetupQuestion
}

func NewClient() (ret *Client) {

	label := "Jina AI"

	ret = &Client{
		PluginBase: &plugins.PluginBase{
			Name:             label,
			SetupDescription: "Jina AI Service - to grab a webpage as clean, LLM-friendly text",
			EnvNamePrefix:    plugins.BuildEnvVariablePrefix(label),
		},
	}

	ret.ApiKey = ret.AddSetupQuestion("API Key", false)

	return
}

// ScrapeURL return the main content of a webpage in clean, LLM-friendly text.
func (jc *Client) ScrapeURL(url string) (ret string, err error) {
	return jc.request(fmt.Sprintf("https://r.jina.ai/%s", url))
}

func (jc *Client) ScrapeQuestion(question string) (ret string, err error) {
	return jc.request(fmt.Sprintf("https://s.jina.ai/%s", question))
}

func (jc *Client) request(requestURL string) (ret string, err error) {
	// Use the client's default API key
	return requestWithApiKey(requestURL, jc.ApiKey.Value)
}

// ScrapeURLWithApiKey scrapes a URL using a custom API key (for delegated execution)
func ScrapeURLWithApiKey(url string, apiKey string) (ret string, err error) {
	return requestWithApiKey(fmt.Sprintf("https://r.jina.ai/%s", url), apiKey)
}

// ScrapeQuestionWithApiKey searches the web using a custom API key (for delegated execution)
func ScrapeQuestionWithApiKey(question string, apiKey string) (ret string, err error) {
	return requestWithApiKey(fmt.Sprintf("https://s.jina.ai/%s", question), apiKey)
}

// requestWithApiKey makes a request to Jina AI with a specific API key
func requestWithApiKey(requestURL string, apiKey string) (ret string, err error) {
	var req *http.Request
	if req, err = http.NewRequest("GET", requestURL, nil); err != nil {
		err = fmt.Errorf("error creating request: %w", err)
		return
	}

	// if api keys exist, set the header
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{}
	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		err = fmt.Errorf("error sending request: %w", err)
		return
	}
	defer resp.Body.Close()

	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		err = fmt.Errorf("error reading response body: %w", err)
		return
	}
	ret = string(body)
	return
}
