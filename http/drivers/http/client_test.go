package http

import (
	"github.com/ArtisanCloud/PowerLibs/v3/http/contract"
	"testing"
)

func Test_NewClient(t *testing.T) {
	helper, err := NewHttpClient(&contract.ClientConfig{})
	if err != nil {
		t.Error(err)
	}

	if helper == nil {
		t.Error(err)
	}

}
