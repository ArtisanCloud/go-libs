package dataflow

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/ArtisanCloud/PowerLibs/v3/http/contract"
	"github.com/pkg/errors"
)

type Dataflow struct {
	client           contract.ClientInterface
	middlewareHandle contract.RequestMiddleware
	request          *http.Request
	option           *Option
	err              []error
}

type Option struct {
	BaseUrl string
}

func NewDataflow(client contract.ClientInterface, middlewareHandle contract.RequestMiddleware, option *Option) *Dataflow {
	if middlewareHandle == nil {
		middlewareHandle = func(handle contract.RequestHandle) contract.RequestHandle {
			return handle
		}
	}

	df := Dataflow{
		client:           client,
		middlewareHandle: middlewareHandle,
		request: &http.Request{
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
		},
		option: option,
	}
	if option == nil {
		return &df
	}
	if option.BaseUrl != "" {
		u, err := url.ParseRequestURI(option.BaseUrl)
		if err != nil {
			df.err = append(df.err, errors.Wrap(err, "base url invalid"))
		}
		df.request.URL = u
		df.request.Host = u.Host
	}
	return &df
}

func (d *Dataflow) getClient() contract.ClientInterface {
	return d.client
}

func (d *Dataflow) getMiddlewareHandle() contract.RequestMiddleware {
	return d.middlewareHandle
}

func (d *Dataflow) WithContext(ctx context.Context) contract.RequestDataflowInterface {
	d.request = d.request.WithContext(ctx)
	return d
}

func (d *Dataflow) Method(method string) contract.RequestDataflowInterface {
	d.request.Method = method
	return d
}

// Uri 请注意 Url 与 Uri 方法是冲突的, Uri方法将 Uri 拼接在 BaseUrl 之后
func (d *Dataflow) Uri(uri string) contract.RequestDataflowInterface {
	if d.option.BaseUrl != "" {
		u, _ := url.ParseRequestURI(d.option.BaseUrl)
		d.request.URL = u
	}
	if d.request.URL == nil {
		d.err = append(d.err, errors.New("invalid request url"))
		return d
	}
	newUrl, err := d.request.URL.Parse(uri)
	if err != nil {
		d.err = append(d.err, err)
		return d
	}
	d.request.URL = newUrl
	return d
}

func (d *Dataflow) Url(requestUrl string) contract.RequestDataflowInterface {
	u, err := url.ParseRequestURI(requestUrl)
	if err != nil {
		d.err = append(d.err, errors.Wrap(err, "invalid url"))
		return d
	}
	d.request.URL = u
	d.request.Host = u.Host
	return d
}

func (d *Dataflow) makeHeaderIfNil() {
	if d.request.Header == nil {
		d.request.Header = make(http.Header)
	}
}

// Header 设置请求头, 对一个 Key 多次调用该方法, values 始终会被后面调用的覆盖
func (d *Dataflow) Header(key string, values ...string) contract.RequestDataflowInterface {
	if len(values) == 0 {
		return d
	}
	d.makeHeaderIfNil()
	for i, v := range values {
		if i == 0 {
			d.request.Header.Set(key, v)
		} else {
			d.request.Header.Add(key, v)
		}
	}
	return d
}

func (d *Dataflow) Query(key string, values ...string) contract.RequestDataflowInterface {
	if len(values) == 0 {
		return d
	}
	query := d.request.URL.Query()
	for i, v := range values {
		if i == 0 {
			query.Set(key, v)
		} else {
			query.Add(key, v)
		}
	}
	d.request.URL.RawQuery = query.Encode()
	return d
}

func (d *Dataflow) Json(jsonAny interface{}) contract.RequestDataflowInterface {
	// 设置 Header
	d.Header("content-type", "application/json")
	d.Header("Accept", "*/*")

	// 标准库Json编码 body reader
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(jsonAny); err != nil {
		d.err = append(d.err, errors.Wrap(err, "json body encode failed"))
		return d
	}
	d.Body(&buf)
	return d
}

func (d *Dataflow) Body(body io.Reader) contract.RequestDataflowInterface {
	if body != nil {
		d.request.Body = io.NopCloser(body)

		switch v := body.(type) {
		case *bytes.Buffer:
			d.request.ContentLength = int64(v.Len())
			buf := v.Bytes()
			d.request.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return io.NopCloser(r), nil
			}
		case *bytes.Reader:
			d.request.ContentLength = int64(v.Len())
			snapshot := *v
			d.request.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return io.NopCloser(&r), nil
			}
		case *strings.Reader:
			d.request.ContentLength = int64(v.Len())
			snapshot := *v
			d.request.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return io.NopCloser(&r), nil
			}
		default:
		}

		if d.request.GetBody != nil && d.request.ContentLength == 0 {
			d.request.Body = http.NoBody
			d.request.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
		}
	}
	return d
}

func (d *Dataflow) Any(data contract.BodyEncoder) contract.RequestDataflowInterface {
	body, err := data.Encode()
	if err != nil {
		d.err = append(d.err, errors.Wrap(err, "body encode failed"))
	}
	d.Body(body)
	return d
}

func (d *Dataflow) Xml(xmlAny interface{}) contract.RequestDataflowInterface {
	// 设置 Header
	d.Header("content-type", "application/xml")

	var buf bytes.Buffer
	encoder := xml.NewEncoder(&buf)
	err := encoder.Encode(xmlAny)
	if err != nil {
		d.err = append(d.err, err)
	}
	d.Body(&buf)
	return d
}

func (d *Dataflow) Multipart(multipartDf func(multipart contract.MultipartDfInterface)) contract.RequestDataflowInterface {
	multipart := NewMultipartHelper()
	multipartDf(multipart)
	err := multipart.Close()
	if err != nil {
		d.err = append(d.err, err)
		return d
	}
	err = multipart.Err()
	if err != nil {
		d.err = append(d.err, err)
		return d
	}
	d.Header("content-type", multipart.GetContentType())
	d.Body(multipart.GetReader())
	return d
}

func (d *Dataflow) Err() error {
	if len(d.err) > 0 {
		return d.err[0]
	}
	return nil
}

func (d *Dataflow) Request() (response *http.Response, err error) {
	if d.Err() != nil {
		return nil, d.Err()
	}

	d.Header("Accept", "*/*")

	handle := d.middlewareHandle(func(request *http.Request) (response *http.Response, err error) {
		return d.client.DoRequest(request)
	})
	resp, err := handle(d.request)
	if err != nil {
		d.err = append(d.err, errors.Wrap(err, "request failed"))
		return resp, d.Err()
	}
	return resp, nil
}

// Result 实现了 Json 解码
func (d *Dataflow) Result(result interface{}) (err error) {
	if result == nil {
		return errors.New("nil result")
	}
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Ptr {
		return errors.New("result is not pointer")
	}
	// request
	resp, err := d.Request()
	if err != nil {
		return err
	}

	// decode 不支持 array
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(result)
	if err != nil {
		return errors.Wrap(err, "decode response failed")
	}
	return nil
}

type Response struct {
	res *http.Response
}

func (r *Response) GetStatusCode() int {
	return r.res.StatusCode
}

func (r *Response) GetHeader(key string) string {
	return r.res.Header.Get(key)
}

func (r *Response) GetBody() io.Reader {
	return r.res.Body
}

func (r *Response) GetBodyBytes() ([]byte, error) {
	if r.res.Body == nil {
		return nil, nil
	}
	body, err := io.ReadAll(r.res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read body failed")
	}
	return body, nil
}

func (r *Response) GetBodyJsonAsMap() (map[string]interface{}, error) {
	data := make(map[string]interface{})
	if r.res.Body == nil {
		return data, nil
	}
	decoder := json.NewDecoder(r.res.Body)
	err := decoder.Decode(&data)
	if err != nil {
		return nil, errors.Wrap(err, "decode body failed")
	}
	return data, nil
}

func (d *Dataflow) RequestResHelper() (response contract.ResponseHelper, err error) {
	resp, err := d.Request()
	if err != nil {
		return nil, err
	}
	return &Response{
		res: resp,
	}, nil
}
