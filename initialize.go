package main

import (
	"encoding/json"
	"fmt"
    "io/ioutil"
	"net"
	"os"
	"strconv"
)
// Things that get called in init() go here

type Config struct {
	ServerConfig struct {
		Address string `json:"Address"`
		Port    string `json:"Port"`
	} `json:"ServerConfig"`
}

//type ServerConfig struct {
//    Address string `json:"Address"`
//    Port    string `json:"Port"`
//} `json:"ServerConfig"`

var config Config

func loadConfigFromFile(f string) {
    if f == "" {
        f = "./config.json"
    }
    configFile, err := ioutil.ReadFile(f)
    if err != nil {
        fmt.Println(os.Stderr, err)
        os.Exit(1)
    }
      
   json.Unmarshal(configFile, &config)
      
   parseAndCheckListenIP(config.ServerConfig.Address)
   parseAndCheckPort(config.ServerConfig.Port)
}

// Pass a string here. Checks IP addresses bound on
// host, and errors if provided IP is not in that
// list
func parseAndCheckListenIP(ip string) {
	//
	addr = net.ParseIP(ip)
	configIpOnHost := false
	// Loop through IPs bound on this system:
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}
	if addr.String() == "0.0.0.0" {
		fmt.Println("Listening on all addresses on system\n")
	} else {
		for _, i := range ifaces {
			ips, err := i.Addrs()
			if err != nil {
				fmt.Println(os.Stderr, err)
				os.Exit(1)
			}
			for _, ipv := range ips {
				var i net.IP
				switch v := ipv.(type) {
				case *net.IPNet:
					i = v.IP
				case *net.IPAddr:
					i = v.IP
				}
				if addr.Equal(i) {
					fmt.Println("Listening on ", i.String(), "\n")
					configIpOnHost = true
					break
				}
			}
		}
	}

	if configIpOnHost == false {
		fmt.Println(os.Stderr, "Configured host IP not bound on host\n")
		os.Exit(1)
	}
}

func parseAndCheckPort(p string) {
	portIsValid := false
	portnum, err := strconv.Atoi(p)
	if err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}

	// If port requires privileges, let's make sure we have them
	if portnum >= 1 && portnum <= 1024 {
		if os.Getuid() == 0 {
			portIsValid = true
		} else {
			fmt.Println(os.Stderr, "Can't bind to port ", p, ". Must be run with root privileges.")
			os.Exit(1)
		}
		// Otherwise, just make sure we're in range of valid ports
	} else if portnum > 1024 && portnum <= 32768 {
		portIsValid = true
	}

	if portIsValid == false {
		fmt.Println(os.Stderr, "Port invalid, must be between 1 and 32768")
		os.Exit(1)
	} else {
		port = p
	}
}