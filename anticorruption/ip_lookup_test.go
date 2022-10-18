package anticorruption

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func Test_IP(t *testing.T) {
	cli := NewIpApi("./ip2region.xdb")
	rsp, err := cli.Query("20.205.243.168")
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(rsp)
}

func Test_a(t *testing.T)  {
	s, err := net.ResolveUDPAddr("udp","127.0.0.1:8989")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(s.IP.String())
}

func GetLANDevIP(devName string) ([]string, error) {
	cmd := exec.Command("ip", "neigh", "flush", "dev", devName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	time.Sleep(time.Second * time.Duration(5))
	//目前openwrt上似乎不是安装的arp命令，而是软链接的查看该文件，所以arp命令实际上应该是如下内容
	cmd = exec.Command("cat", "/proc/1/net/arp")
	output, err = cmd.Output()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(string(output))

	var res []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		if strings.Contains(scanner.Text(), devName) {
			tmp := strings.Split(scanner.Text(), " ")
			res = append(res, tmp[0])
		}
	}
	return res, nil
}