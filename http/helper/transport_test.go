package helper

import (
	"fmt"
	"github.com/ArtisanCloud/PowerLibs/v3/http/contract"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type MyTransport struct {
	http.Transport
}

func (t *MyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	trip, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	// 这两行的意思就是，只要发起了任何请求，我这里都能获取到请求参数和响应内容
	// 这样就能给 otel 做监控了
	fmt.Println("request ++++", req.URL)
	fmt.Println("resp status", trip.Status)
	return trip, err
}

func TestName(t *testing.T) {
	requestHelper, err := NewRequestHelper(&Config{
		BaseUrl: "http://baidu.com",
		ClientConfig: &contract.ClientConfig{
			Transport: &MyTransport{},
		},
	})
	assert.NoError(t, err)
	request, err := requestHelper.Df().Request()
	assert.NoError(t, err)
	fmt.Println(request)
}
