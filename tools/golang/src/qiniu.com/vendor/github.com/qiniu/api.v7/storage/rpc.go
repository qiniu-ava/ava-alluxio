package storage

// The original library rpc.v7 logic in github.com/qiniu/x has its own bugs
// under the concurrent http calls, we make a fork of the library and fix
// the bug
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/conf"
	"github.com/qiniu/x/reqid.v7"
	. "golang.org/x/net/context"
)

var UserAgent = "Golang qiniu/rpc package"
var DefaultClient = Client{&http.Client{Transport: http.DefaultTransport}}

// --------------------------------------------------------------------

type Client struct {
	*http.Client
}

// userApp should be [A-Za-z0-9_\ \-\.]*
func SetAppName(userApp string) error {
	UserAgent = fmt.Sprintf(
		"QiniuGo/%s (%s; %s; %s) %s", conf.Version, runtime.GOOS, runtime.GOARCH, userApp, runtime.Version())
	return nil
}

// --------------------------------------------------------------------

func newRequest(ctx Context, method, reqUrl string, headers http.Header, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, reqUrl, body)
	if err != nil {
		return
	}

	if headers == nil {
		headers = http.Header{}
	}

	req.Header = headers

	//check access token
	mac, ok := ctx.Value("mac").(*qbox.Mac)
	if ok {
		token, signErr := mac.SignRequest(req)
		if signErr != nil {
			err = signErr
			return
		}
		req.Header.Add("Authorization", "QBox "+token)
	}

	return
}

func (r Client) DoRequest(ctx Context, method, reqUrl string, headers http.Header) (resp *http.Response, err error) {
	req, err := newRequest(ctx, method, reqUrl, headers, nil)
	if err != nil {
		return
	}
	return r.Do(ctx, req)
}

func (r Client) DoRequestWith(ctx Context, method, reqUrl string, headers http.Header, body io.Reader,
	bodyLength int) (resp *http.Response, err error) {

	req, err := newRequest(ctx, method, reqUrl, headers, body)
	if err != nil {
		return
	}
	req.ContentLength = int64(bodyLength)
	return r.Do(ctx, req)
}

func (r Client) DoRequestWith64(ctx Context, method, reqUrl string, headers http.Header, body io.Reader,
	bodyLength int64) (resp *http.Response, err error) {

	req, err := newRequest(ctx, method, reqUrl, headers, body)
	if err != nil {
		return
	}
	req.ContentLength = bodyLength
	return r.Do(ctx, req)
}

func (r Client) DoRequestWithForm(ctx Context, method, reqUrl string, headers http.Header,
	data map[string][]string) (resp *http.Response, err error) {

	if headers == nil {
		headers = http.Header{}
	}
	headers.Add("Content-Type", "application/x-www-form-urlencoded")

	requestData := url.Values(data).Encode()
	if method == "GET" || method == "HEAD" || method == "DELETE" {
		if strings.ContainsRune(reqUrl, '?') {
			reqUrl += "&"
		} else {
			reqUrl += "?"
		}
		return r.DoRequest(ctx, method, reqUrl+requestData, headers)
	}

	return r.DoRequestWith(ctx, method, reqUrl, headers, strings.NewReader(requestData), len(requestData))
}

func (r Client) DoRequestWithJson(ctx Context, method, reqUrl string, headers http.Header,
	data interface{}) (resp *http.Response, err error) {

	reqBody, err := json.Marshal(data)
	if err != nil {
		return
	}

	if headers == nil {
		headers = http.Header{}
	}
	headers.Add("Content-Type", "application/json")
	return r.DoRequestWith(ctx, method, reqUrl, headers, bytes.NewReader(reqBody), len(reqBody))
}

func (r Client) Do(ctx Context, req *http.Request) (resp *http.Response, err error) {

	if ctx == nil {
		ctx = Background()
	}

	if reqId, ok := reqid.FromContext(ctx); ok {
		req.Header.Set("X-Reqid", reqId)
	}

	if _, ok := req.Header["User-Agent"]; !ok {
		req.Header.Set("User-Agent", UserAgent)
	}

	transport := r.Transport // don't change r.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	// avoid cancel() is called before Do(req), but isn't accurate
	select {
	case <-ctx.Done():
		err = ctx.Err()
		return
	default:
	}

	if tr, ok := getRequestCanceler(transport); ok {
		// support CancelRequest
		reqC := make(chan bool, 1)
		go func() {
			resp, err = r.Client.Do(req)
			reqC <- true
		}()
		select {
		case <-reqC:
		case <-ctx.Done():
			tr.CancelRequest(req)
			<-reqC
			err = ctx.Err()
		}
	} else {
		resp, err = r.Client.Do(req)
	}
	return
}

// --------------------------------------------------------------------

type ErrorInfo struct {
	Err   string `json:"error,omitempty"`
	Key   string `json:"key,omitempty"`
	Reqid string `json:"reqid,omitempty"`
	Errno int    `json:"errno,omitempty"`
	Code  int    `json:"code"`
}

func (r *ErrorInfo) ErrorDetail() string {

	msg, _ := json.Marshal(r)
	return string(msg)
}

func (r *ErrorInfo) Error() string {

	return r.Err
}

func (r *ErrorInfo) RpcError() (code, errno int, key, err string) {

	return r.Code, r.Errno, r.Key, r.Err
}

func (r *ErrorInfo) HttpCode() int {

	return r.Code
}

// --------------------------------------------------------------------

func parseError(e *ErrorInfo, r io.Reader) {

	body, err1 := ioutil.ReadAll(r)
	if err1 != nil {
		e.Err = err1.Error()
		return
	}

	var ret struct {
		Err   string `json:"error"`
		Key   string `json:"key"`
		Errno int    `json:"errno"`
	}
	if json.Unmarshal(body, &ret) == nil && ret.Err != "" {
		// qiniu error msg style returns here
		e.Err, e.Key, e.Errno = ret.Err, ret.Key, ret.Errno
		return
	}
	e.Err = string(body)
}

func ResponseError(resp *http.Response) (err error) {

	e := &ErrorInfo{
		Reqid: resp.Header.Get("X-Reqid"),
		Code:  resp.StatusCode,
	}
	if resp.StatusCode > 299 {
		if resp.ContentLength != 0 {
			ct, ok := resp.Header["Content-Type"]
			if ok && strings.HasPrefix(ct[0], "application/json") {
				parseError(e, resp.Body)
			}
		}
	}
	return e
}

func CallRet(ctx Context, ret interface{}, resp *http.Response) (err error) {

	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode/100 == 2 {
		// fmt.Println("CallRet status code is 2xx", resp.Request.URL.String(), resp)
		b := make([]byte, 10000, 100000)
		i, e := resp.Body.Read(b)
		if e == nil {
			fmt.Println(i, e, string(b[:i]))
		} else {
			if i > 0 {
				fmt.Println(e, string(b[:i]))
			}
		}

		if ret != nil && resp.ContentLength != 0 {
			fmt.Println("I'm here")
			err = json.Unmarshal(b[:i], ret)
			fmt.Printf("I'm here again with ret: %v\n", ret)
			// err = json.NewDecoder(resp.Body).Decode(ret)
			if err != nil {
				return
			}
		}
		if resp.StatusCode == 200 {
			return nil
		}
	}
	return ResponseError(resp)
}

func (r Client) CallWithForm(ctx Context, ret interface{}, method, reqUrl string, headers http.Header,
	param map[string][]string) (err error) {

	resp, err := r.DoRequestWithForm(ctx, method, reqUrl, headers, param)
	if err != nil {
		return err
	}
	return CallRet(ctx, ret, resp)
}

func (r Client) CallWithJson(ctx Context, ret interface{}, method, reqUrl string, headers http.Header,
	param interface{}) (err error) {

	resp, err := r.DoRequestWithJson(ctx, method, reqUrl, headers, param)
	if err != nil {
		return err
	}
	return CallRet(ctx, ret, resp)
}

func (r Client) CallWith(ctx Context, ret interface{}, method, reqUrl string, headers http.Header, body io.Reader,
	bodyLength int) (err error) {

	resp, err := r.DoRequestWith(ctx, method, reqUrl, headers, body, bodyLength)
	if err != nil {
		return err
	}
	return CallRet(ctx, ret, resp)
}

func (r Client) CallWith64(ctx Context, ret interface{}, method, reqUrl string, headers http.Header, body io.Reader,
	bodyLength int64) (err error) {

	resp, err := r.DoRequestWith64(ctx, method, reqUrl, headers, body, bodyLength)
	if err != nil {
		return err
	}
	return CallRet(ctx, ret, resp)
}

func (r Client) Call(ctx Context, ret interface{}, method, reqUrl string, headers http.Header) (err error) {

	resp, err := r.DoRequestWith(ctx, method, reqUrl, headers, nil, 0)
	if err != nil {
		return err
	}
	return CallRet(ctx, ret, resp)
}

// ---------------------------------------------------------------------------

type requestCanceler interface {
	CancelRequest(req *http.Request)
}

type nestedObjectGetter interface {
	NestedObject() interface{}
}

func getRequestCanceler(tp http.RoundTripper) (rc requestCanceler, ok bool) {

	if rc, ok = tp.(requestCanceler); ok {
		return
	}

	p := interface{}(tp)
	for {
		getter, ok1 := p.(nestedObjectGetter)
		if !ok1 {
			return
		}
		p = getter.NestedObject()
		if rc, ok = p.(requestCanceler); ok {
			return
		}
	}
}

// --------------------------------------------------------------------
