package api

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type KDNiaoConfig struct {
	KdnBusinessId string `split_words:"true" required:"true"`
	KdnApiSecret  string `split_words:"true" required:"true"`
}

// KDNNiaoApi provides some methods to get parcel info from
// 快递鸟
type KDNiaoApi struct {
	KDNiaoConfig
	httpC *http.Client
}

func NewKDNiaoApi(config KDNiaoConfig) *KDNiaoApi {
	return &KDNiaoApi{
		KDNiaoConfig: config,
		httpC: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// SignedRequest returns hashed API secret and data to be used as authentication token
func (kd KDNiaoApi) SignedRequest(req []byte) string {
	str := string(req) + kd.KdnApiSecret
	hash := md5.Sum([]byte(str))

	signedData := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%x", hash)))
	return signedData
}

// GetSourceByTrack fetches information of the parcel by its code
func (kd KDNiaoApi) GetSourceByTrack(code string) (*SourceResponse, error) {
	reqData := map[string]string{
		"LogisticCode": code,
	}

	jsonReq, _ := json.Marshal(reqData)
	reqStr := "http://api.kdniao.com/Ebusiness/EbusinessOrderHandle.aspx?EBusinessID=%s&DataType=%d&DataSign=%s&RequestType=%d&RequestData=%s"
	reqUrl := fmt.Sprintf(reqStr, kd.KdnBusinessId, 2, url.QueryEscape(kd.SignedRequest(jsonReq)), 2002, url.QueryEscape(string(jsonReq)))

	u, err := url.Parse(reqUrl)
	if err != nil {
		return nil, errors.WithMessage(err, "wrong url")
	}

	req, _ := http.NewRequest("POST", u.String(), nil)

	resp, err := kd.httpC.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsonResp SourceResponse

	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp, nil
}

type SourceResponse struct {
	LogisticCode string    `json:"LogisticCode"`
	Shippers     []Shipper `json:"Shippers"`
	EBusinessID  string    `json:"EBusinessID"`
	Code         string    `json:"Code"`
	Success      bool      `json:"Success"`
}

type Shipper struct {
	ShipperName string `json:"ShipperName"`
	ShipperCode string `json:"ShipperCode"`
}
