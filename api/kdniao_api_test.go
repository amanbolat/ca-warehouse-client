package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestKDNiaoApi_GetSourceByTrack(t *testing.T) {
	api := KDNiaoApi{
		BusinessId: "1613169",
		ApiSecret:  "11f4d1b8-11cd-44da-9ba9-5ae2a32de09d",
		httpC: &http.Client{
			Timeout: time.Second * 5,
		},
	}

	resp, err := api.GetSourceByTrack("SF1192131829831")
	assert.NoError(t, err)
	t.Log(resp)
}
