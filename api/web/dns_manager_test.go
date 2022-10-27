package web

import (
	"encoding/json"
	"fmt"
	"github.com/hhzhhzhhz/local-dns/utils"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func Test_DnsA(t *testing.T)  {
}

func do(uri string, r io.Reader) (string, error) {
	req, err := http.NewRequest("POST", uri, r)
	if err != nil {
		return "", err
	}
	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		resp.Body.Close()
	}()
	b, err := ioutil.ReadAll(resp.Body)
	rsp := utils.Response{}
	if err := json.Unmarshal(b, &rsp); err != nil {
		return "", err
	}
	return fmt.Sprintf("code:%d info:%s, data:%v", rsp.Code, rsp.Info, rsp.Data), nil
}

type DnsReq struct {
	Key string `json:"key"`
	Host string `json:"host"`
	Text string `json:"text"`
}

func Test_RemoveDNS(t *testing.T) {
	uri := "http://127.0.0.1:9876/api/dnsMgr/RemoveDnsRecord?domains=['aaa.com']"
	rs, err := do(uri, strings.NewReader(string("")))
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(rs)

}