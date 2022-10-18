package anticorruption

import (
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"sync"
)

var ipApi *IpApi

func GetIpApi() *IpApi {
	return ipApi
}

func Init(dbPath string)  {
	ipApi = NewIpApi(dbPath)
}

type IpInfo struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
	ISP     string `json:"isp"`
}

type IpApi struct {
	mux sync.Mutex
	dbPath string
	s *xdb.Searcher
}

func NewIpApi(dbPath string) *IpApi {
	return &IpApi{
		dbPath: dbPath,
	}
}

func (i *IpApi) Query(ip string) (result string, err error) {
	i.mux.Lock()
	defer i.mux.Unlock()
	s, err := i.cli()
	if err != nil {
		return result, err
	}
	result, err = s.SearchByStr(ip)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (i *IpApi) cli() (*xdb.Searcher, error) {
	if i.s != nil {
		return i.s, nil
	}
	buf, err := xdb.LoadContentFromFile(i.dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load content from `%s`: %s", i.dbPath, err)
	}
	s, err := xdb.NewWithBuffer(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create searcher with content: %s", err)
	}
	i.s = s
	return i.s, nil
}

