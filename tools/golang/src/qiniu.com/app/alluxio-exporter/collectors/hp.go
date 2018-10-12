package collectors

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

func mapToURLValues(m map[string]string) (u url.Values) {
	u = url.Values{}
	if m == nil || len(m) == 0 {
		return nil
	}
	for key, value := range m {
		u.Set(key, value)
	}
	return
}

func HTTPRequest(URL, method string, data, query map[string]string) (body []byte, e error) {

	queryValues := mapToURLValues(query)
	URL += queryValues.Encode()

	req, e := http.NewRequest(method, URL, nil)
	if e != nil {
		log.Println("[WARN] cannot make new request:", e)
		return nil, e
	}
	req.PostForm = mapToURLValues(data)
	req.ParseForm()

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	resp, e := client.Do(req)
	if e != nil {
		log.Println("[WARN] cannot receive response:", e)
		return nil, e
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		log.Println("[WARN] response status unsuccess:", resp.StatusCode)
		e = errors.Errorf("Unsuccess Response Code: " + strconv.Itoa(resp.StatusCode))
		return nil, e
	}

	body, e = ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Println("[WARN] cannot read body:", e)
		return nil, e
	}

	return
}
