package calendly

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Calendly struct {
	APIKey string
}

func NewCalendly() *Calendly {
	return &Calendly{}
}

func (c *Calendly) SetAPIKey(key string) {
	c.APIKey = key
}

func (c *Calendly) SetAPIKeyFromFile(fn string) error {
	key, err := readKeyFromFile(fn)
	if err != nil {
		return err
	}
	c.APIKey = key
	return nil
}

func (c *Calendly) GetMe() (map[string]interface{}, error) {
	me_raw, err := c.APIRequest("users/me", nil)
	if err != nil {
		return nil, err
	}

	var me interface{}
	err = json.Unmarshal(me_raw, &me)
	if err != nil {
		return nil, err
	}

	return GetStringMap(me, "resource"), nil
}

func (c *Calendly) GetEvents() ([]*Event, error) {
	me, err := c.GetMe()
	if err != nil {
		return nil, err
	}

	if _, ok := me["uri"]; !ok {
		return nil, fmt.Errorf("cannot determine user uri")
	}

	params := map[string]interface{}{
		"count":  "100",
		"status": "active",
		"user":   me["uri"],
	}
	j, err := c.APIRequestToJSON("scheduled_events", params)
	if err != nil {
		return nil, err
	}
	for x := range j {
		fmt.Println(x)
	}
	if c, ok := j["collection"]; ok {
		events, err := populateEventsFromJSON(c)
		if err != nil {
			return nil, err
		}
		return events, nil
	}

	return nil, fmt.Errorf("no collection in scheduled_events response")
}

func (c *Calendly) APIRequestToJSON(method string, params map[string]interface{}) (map[string]interface{}, error) {
	raw, err := c.APIRequest(method, params)
	if err != nil {
		return nil, err
	}

	var j interface{}
	err = json.Unmarshal(raw, &j)
	if err != nil {
		return nil, err
	}

	return j.(map[string]interface{}), nil
}

func (c *Calendly) APIRequest(method string, params map[string]interface{}) ([]byte, error) {
	url := "https://api.calendly.com/" + method
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if c.APIKey != "" {
		req.Header.Add("Authorization", "Bearer "+c.APIKey)
	}
	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v.(string))
		}
		req.URL.RawQuery = q.Encode()
	}
	rep, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if rep.StatusCode != 200 {
		return nil, fmt.Errorf("error %v requesting %v", rep.StatusCode, req.URL)
	}
	repBody, err := io.ReadAll(rep.Body)
	if err != nil {
		return nil, err
	}
	return repBody, nil
}

func readKeyFromFile(fn string) (string, error) {
	key, err := os.ReadFile(fn)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(key), "\r\n"), err
}
