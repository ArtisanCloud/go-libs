package http

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/ArtisanCloud/PowerLibs/v3/fmt"
	"github.com/ArtisanCloud/PowerLibs/v3/http/contract"
	"github.com/pkg/errors"
)

// Client 是 net/http 的封装
type Client struct {
	conf       *contract.ClientConfig
	coreClient *http.Client
}

func NewHttpClient(config *contract.ClientConfig) (*Client, error) {
	if config == nil {
		config = &contract.ClientConfig{}
		config.Default()
	}
	coreClient := http.Client{
		Timeout: config.Timeout,
	}
	var proxy func(*http.Request) (*url.URL, error)
	if config.ProxyURI != "" {
		if proxyURL, err := url.Parse(config.ProxyURI); err == nil {
			proxy = http.ProxyURL(proxyURL)
		}
	}
	if config.Cert.CertFile != "" && config.Cert.KeyFile != "" {
		certPair, err := tls.LoadX509KeyPair(config.Cert.CertFile, config.Cert.KeyFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load certificate")
		}
		coreClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{certPair},
			}, Proxy: proxy}
	} else if proxy != nil {
		coreClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}, Proxy: proxy}
	}

	if config.Transport != nil {
		// 这里会直接覆盖上面的 tls 和 proxy 有能力配置这个的话会自己解决的
		coreClient.Transport = config.Transport
	}

	return &Client{
		conf:       config,
		coreClient: &coreClient,
	}, nil
}

// SetConfig 配置客户端
func (c *Client) SetConfig(config *contract.ClientConfig) {
	if config != nil {
		c.conf = config
	}
	var proxy func(*http.Request) (*url.URL, error)
	if config.ProxyURI != "" {
		if proxyURL, err := url.Parse(config.ProxyURI); err == nil {
			proxy = http.ProxyURL(proxyURL)
		}
	}
	coreClient := http.Client{
		Timeout: config.Timeout,
	}
	// todo set coreClient
	if config.Cert.CertFile != "" && config.Cert.KeyFile != "" {
		certPair, err := tls.LoadX509KeyPair(config.Cert.CertFile, config.Cert.KeyFile)
		if err != nil {
			err = errors.Wrap(err, "failed to load certificate")
			fmt.Dump(err)
			return
		}
		coreClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{certPair},
		}}
		c.coreClient = &coreClient
	} else if proxy != nil {
		coreClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}, Proxy: proxy}
	}
}

// GetConfig 返回配置副本
func (c *Client) GetConfig() contract.ClientConfig {
	return *c.conf
}

func (c *Client) DoRequest(request *http.Request) (response *http.Response, err error) {
	return c.coreClient.Do(request)
}
