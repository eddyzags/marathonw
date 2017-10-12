package marathon

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func (m *marathon) apiCall(method, path string, reader io.Reader, result interface{}) error {
	req, err := m.makeRequest(method, path, reader)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return parseError(resp)
	}

	if result != nil {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(b, result); err != nil {
			return err
		}
	}

	return nil
}

func (m *marathon) makeRequest(method, path string, reader io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", m.config.URI, path)

	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	if m.config.HTTPBasicAuthUser != "" && m.config.HTTPBasicAuthPassword != "" {
		request.SetBasicAuth(m.config.HTTPBasicAuthUser, m.config.HTTPBasicAuthPassword)
	}

	if m.config.DCOSToken != "" {
		request.Header.Add("Authorization", "token="+m.config.HTTPBasicAuthUser)
	}

	request.Header.Add("Content-Type", "application/json")

	return request, nil
}
