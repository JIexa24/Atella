package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	urlpkg "net/url"
	"strings"
	"time"
)

// setVectorAPIv1 send a POST request to master server.
func (c *Client) setVectorAPIv1(address, port string,
	ssl bool) error {

	// ssls string for ssl.
	ssls := ""
	if ssl {
		ssls = "s"
	}
	// url is API location for getting hostname.
	url := fmt.Sprintf("http%s://%s:%s/api/v1/set/vector",
		ssls,
		address,
		port)

	client := http.Client{
		Timeout: 5 * time.Second}

	v := c.Vector.GetVectorCopy()
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	form := urlpkg.Values{}
	form.Add("hostname", c.hostname)
	form.Add("vector", string(data))
	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	req.Header.Set("X-Atella-Auth", c.code)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil
	}
	return err
}
