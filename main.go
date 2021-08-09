package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {

	var (
		per             operation
		proxyList4file  []proxyer
		canUseProxyList []proxyer
		finalout        []proxyer
		socksTestDataer socks5TestData
	)
	per.INIT()

	logfile, err := os.OpenFile(per.path.logfilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logfile.Close()

	log.SetOutput(logfile)
	log.Println("[INFO] Programer start.")

	proxyList4file = formatProxylist(per.path.fullPath4ProxyListFile)

	canUseProxyList = firstTestList(proxyList4file, per.timer.tcpTimeout)

	socksTestDataer = socks5TestData{canUseProxyList, per}

	socksTestDataer.socks5Test()

	// finalout = append(finalout, per.staticProxy...)

	for _, v := range socksTestDataer.proxylist {
		if v.allTestok {
			finalout = append(finalout, v)
		}
	}

	switch len(finalout) {
	case 0:
		log.Println("[ERROR]", "no valid proxy server")
		return
	case 1:
		per.outputfile(finalout)
	default:
		per.outputfile(sortProxyList(finalout))
	}

	log.Println("[INFO] Programer stop.")

}

func sortProxyList(inlist []proxyer) []proxyer {
	var outlist []proxyer
	var vt1 []float64
	for _, v := range inlist {
		vt1 = append(vt1, v.speed)
	}
	sort.Float64s(vt1)
	for _, v := range vt1 {
		for _, vv := range inlist {
			if vv.speed == v {
				outlist = append(outlist, vv)
			}
		}
	}
	return outlist
}

func (s *socks5TestData) socks5Test() {
	var vg sync.WaitGroup
	for i, v := range s.proxylist {
		vg.Add(1)
		go func(i int, v proxyer) {
			ss := 0
			var ff float64
			for _, vv := range s.pers.testwebslist {
				vt5f := proxytest(v.fullstr, vv.url, s.pers.timer.httpTimeout)
				if vt5f != 0 {
					ss += 1
					ff += vt5f
				}
			}
			if ss == len(s.pers.testwebslist) {
				s.proxylist[i].allTestok = true
				s.proxylist[i].speed = ff
			}
			vg.Done()
		}(i, v)
	}
	vg.Wait()
}

func firstTestList(p1 []proxyer, t1 time.Duration) []proxyer {
	var reproxylist []proxyer
	var vg sync.WaitGroup
	for i, v := range p1 {
		vg.Add(1)
		go func(i int, v proxyer) {
			if v.connectTest(t1) {
				reproxylist = append(reproxylist, v)
			}
			vg.Done()
		}(i, v)
	}
	vg.Wait()
	log.Printf("[INFO] [ %d/%d live] first test complete .\n", len(p1), len(reproxylist))
	return reproxylist
}

func formatProxylist(filename string) []proxyer {
	var reList []proxyer

	f, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return reList
	}
	if bytes.Count(f, []byte("\n")) < 1 {
		log.Println(errors.New("[ERROR]Can not find proxy server."))
		return reList
	}
	vts1 := strings.Split(string(f), "\n")
	for _, v := range vts1 {

		if strings.Index(v, ":") == -1 || v == "" {
			continue
		}

		for _, vv := range []string{" ", "\r", "\n"} {
			v = strings.ReplaceAll(v, vv, "")
		}

		var vtp3 proxyer
		vtp3.socksstr = v
		vtp3.stype = "socks5"
		vti4 := strings.Index(vtp3.socksstr, ":")
		if vti4 != -1 {
			vtp3.ipaddress = vtp3.socksstr[:vti4]
			vtp3.ports = vtp3.socksstr[vti4+1:]
			vtp3.port, _ = strconv.Atoi(vtp3.ports)
		}
		vtp3.fullstr = vtp3.stype + "://" + vtp3.socksstr
		reList = append(reList, vtp3)
	}
	log.Printf("[INFO] [%s] readed %d item,useable %d item.\n", filename, len(vts1), len(reList))
	return reList
}
