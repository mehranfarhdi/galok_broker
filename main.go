//package main
//
//import "github.com/hauke96/kingpin"
//
//const VERSION string = "v0.0.4"
//
//var (
//	app           = kingpin.New("goMS", "Galok Message Broker write with go")
//	appConfigFile = app.Flag("config", "Specifies the configuration file that should be used. This is \"sample.env\" by default.").Short('c').Default("./sample.env").String()
//)
//
//func main() {
//
//}

package main

import (
	"fmt"
	"golang.org/x/net/proxy"
	"net/http"
	"os"
)

func main() {
	socksProxy := "socks5://127.0.0.1:1080" // Replace with your SOCKS5 proxy address and port
	proxyDialer, err := proxy.SOCKS5("tcp", socksProxy, nil, proxy.Direct)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating SOCKS5 proxy dialer: %v\n", err)
		os.Exit(1)
	}

	transport := &http.Transport{Dial: proxyDialer.Dial}
	client := &http.Client{Transport: transport}

	modulePath := "github.com/jinzhu/gorm"
	resp, err := client.Get(fmt.Sprintf("https://proxy.golang.org/%s/@v/list", modulePath))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching module versions: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	fmt.Println("Module versions:")
	// Handle the response body as needed
}
