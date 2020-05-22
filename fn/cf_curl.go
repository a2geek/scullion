package fn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"scullion/ctx"
	"strings"

	"github.com/cloudfoundry-community/go-cfclient"
)

func NewCfCurlRegistrar(client *cfclient.Client) Registrar {
	return func(state *ctx.State) {
		cf := CfCurl{
			state:  state,
			client: client,
		}
		state.AddFunc("GET", cf.Get)
		state.AddFunc("GetResources", cf.GetResources)
		state.AddFunc("POST", cf.Post)
		state.AddFunc("PUT", cf.Put)
	}
}

// CfCurl encompasses all Cloud Foundry interactions.
type CfCurl struct {
	state  *ctx.State
	client *cfclient.Client
}

// Get performs a single HTTP GET against the CF API.
// The entire JSON response is added as 'name' emitted into the
// channel for processing.
func (cf *CfCurl) Get(path, name string) error {
	cf.state.Debugf("%s <- GET %s", name, path)
	doc, err := cf.makeRequest("GET", path)
	if err != nil {
		return err
	}

	cf.state.EmitVar(map[string]interface{}{
		name: doc,
	})

	return nil
}

// Get performs a multi-valued HTTP GET against the CF API inclusive of paging.
// For every item in 'resources[...]' a value of 'name' will be emitted into the
// channel for further processing.
func (cf *CfCurl) GetResources(path, name string) error {
	cf.state.Debugf("%s <- GetResources %s", name, path)
	if path != "" {
		doc, err := cf.makeRequest("GET", path)

		rss, ok := doc["resources"]
		if !ok {
			return fmt.Errorf("GET %s does not contain 'resources'", path)
		}

		ary, ok := rss.([]interface{})
		if !ok {
			return fmt.Errorf("GET %s 'resources' is not an array", path)
		}

		for item := range ary {
			cf.state.EmitVar(map[string]interface{}{
				name: item,
			})
		}

		path, err = cf.nextPage(doc)
		return err
	}
	return nil
}

// Post performs a single HTTP POST against the CF API.
func (cf *CfCurl) Post(path, body string) error {
	cf.state.Debugf("POST %s", path)
	// Note we toss the doc away for now...
	_, err := cf.makeRequestWithBody("POST", path, body)
	return err
}

// Put performs a single HTTP PUT against the CF API.
func (cf *CfCurl) Put(path, body string) error {
	cf.state.Debugf("PUT %s", path)
	// Note we toss the doc away for now...
	_, err := cf.makeRequestWithBody("PUT", path, body)
	return err
}

// Delete performs a single HTTP DELETE against the CF API.
func (cf *CfCurl) Delete(path string) error {
	cf.state.Debugf("DELETE %s", path)
	// Note we toss the doc away for now...
	_, err := cf.makeRequest("DELETE", path)
	return err
}

func (cf *CfCurl) makeRequest(method, path string) (map[string]interface{}, error) {
	return cf.internalMakeRequest(cf.client.NewRequest(method, path))
}

func (cf *CfCurl) makeRequestWithBody(method, path, body string) (map[string]interface{}, error) {
	return cf.internalMakeRequest(cf.client.NewRequestWithBody(method, path, strings.NewReader(body)))
}

func (cf *CfCurl) internalMakeRequest(req *cfclient.Request) (map[string]interface{}, error) {
	resp, err := cf.client.DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("request returned error: %v", err)
	}

	doc, err := cf.readDoc(resp)
	if err != nil {
		return nil, fmt.Errorf("JSON parse returned error: %v", err)
	}

	return doc, err
}

func (cf *CfCurl) readDoc(resp *http.Response) (map[string]interface{}, error) {
	doc := make(map[string]interface{})
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&doc)
	if err != nil {
		return nil, fmt.Errorf("unable to decode JSON document: %v", err)
	}
	return doc, nil
}

// V2 and V3 store this in different places...
func (cf *CfCurl) nextPage(doc map[string]interface{}) (string, error) {
	if path, ok := doc["next_url"]; ok {
		return path.(string), nil
	}

	if page, ok := doc["pagination"]; ok {
		pagination := cfclient.Pagination{}
		err := json.Unmarshal(page.([]byte), &pagination)
		if err != nil {
			return "", err
		}
		if l, ok := pagination.Next.(cfclient.Link); ok {
			return l.Href, nil
		}
		return "", nil
	}

	return "", fmt.Errorf("unable to parse pagination element")
}
