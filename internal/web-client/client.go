package source

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rates_service/common"
)

func GetRates() (common.Depth, error) {
	var depth = common.Depth{}

	data, err := sendRequest(common.SourceAPI)

	if err != nil {
		return depth, err
	}

	err = json.Unmarshal(data, &depth)
	if err != nil {
		return depth, err
	}

	return depth, nil

}

func sendRequest(url string) ([]byte, error) {
	fmt.Println("...sendRequest()")
	defer fmt.Println("...sendRequest() done")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		return data, nil
	}

	return nil, err
}
