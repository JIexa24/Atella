package httpclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// getHostnameAPIv1 request hostname from master server.
func (c *Client) getHostnameAPIv1(address, port string,
	ssl bool) (string, int, error) {
	ssls := ""
	if ssl {
		ssls = "s"
	}
	url := fmt.Sprintf("http%s://%s:%s/api/v1/get/hostname",
		ssls,
		address,
		port)

	client := http.Client{
		Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return "unknown", 0, fmt.Errorf("%s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	response := universalResponse{}
	json.Unmarshal(body, &response)
  hostname := response.ResponseBody.(string)
	return hostname, resp.StatusCode, err
}
