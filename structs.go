package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"time"

	"github.com/spf13/viper"
)

type socks5TestData struct {
	proxylist []proxyer
	pers      operation
}

type proxyer struct {
	stype     string
	ipaddress string
	port      int
	ports     string
	socksstr  string
	fullstr   string
	allTestok bool
	speed     float64
}

type operation struct {
	timer struct {
		tcpTimeout  time.Duration
		httpTimeout time.Duration
	}
	path struct {
		sourceDir              string
		proxyListFile          string
		fullPath4ProxyListFile string
		logfilename            string
		outputFile             bool
		outputDir              string
		outputFilename         string
		outputTitle            string
		outputTail             string
	}
	testwebslist []testwebs
	staticProxy  []proxyer
}

type testwebs struct {
	url    string
	result bool
	speed  float64
}

func (i *operation) INIT() {

	// i.testwebslist = append(i.testwebslist, testwebs{"https://www.baidu.com", false, 0})
	// i.testwebslist = append(i.testwebslist, testwebs{"https://www.google.com", false, 0})
	// i.testwebslist = append(i.testwebslist, testwebs{"https://github.com", false, 0})

	//i.timer.tcpTimeout = 5 * time.Second
	// i.timer.httpTimeout = 10 * time.Second

	// i.path.sourceDir = "./source"
	// i.path.proxyListFile = "p.txt"
	// i.path.fullPath4ProxyListFile = path.Join(i.path.sourceDir, i.path.proxyListFile)

	// i.path.logfilename = "./log/" + time.Now().Format("20060102") + ".log"

	// i.path.outputFile = true
	// i.path.outputDir = "./"
	// i.path.outputFilename = "rc"
	// i.path.outputTitle = "listen = 192.168.19.202:80\n"
	// i.path.outputTail = "#Update by "

	viper.SetConfigName("config.json")
	viper.SetConfigType("json")

	viper.AddConfigPath(".")
	viper.AddConfigPath("/root/.cow")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("[ERROR]", err)
	}

	for _, v := range viper.GetStringSlice("testwebslist") {
		i.testwebslist = append(i.testwebslist, testwebs{v, false, 0})
	}
	for _, v2 := range viper.GetStringSlice("static") {
		i.staticProxy = append(i.staticProxy, proxyer{"", "", 0, "", "", v2, false, 0})
	}

	i.timer.tcpTimeout = time.Duration(viper.GetFloat64("tcpTimeout") * float64(time.Second))
	i.timer.httpTimeout = time.Duration(viper.GetFloat64("httpTimeout") * float64(time.Second))

	i.path.fullPath4ProxyListFile = path.Join(viper.GetString("maindir"), viper.GetString("sourceDir"), viper.GetString("proxyListFile"))
	i.path.outputFile = viper.GetBool("outputFile")
	i.path.logfilename = path.Join(viper.GetString("maindir"), "log", time.Now().Format("20060102")+".log")

	i.path.outputDir = viper.GetString("outputDir")
	i.path.outputFilename = path.Join(i.path.outputDir, viper.GetString("outputFilename"))
	i.path.outputTitle = viper.GetString("outputTitle")
	i.path.outputTail = viper.GetString("outputTail")
	// fmt.Println(i)
}

func (o *operation) outputfile(finalout []proxyer) {
	if o.path.outputFile {
		writetext := o.path.outputTitle

		for _, v := range finalout {
			writetext += fmt.Sprintf("proxy = %s\n# Speed : %.3f s\n", v.fullstr, v.speed)
		}
		writetext += fmt.Sprintf("%s%s\n", o.path.outputTail, time.Now().Format("2006-01-02 15:04:05"))

		err := ioutil.WriteFile(o.path.outputFilename, []byte(writetext), 0644)
		if err != nil {
			log.Println("[ERROR]", err)
		}
	}
}
