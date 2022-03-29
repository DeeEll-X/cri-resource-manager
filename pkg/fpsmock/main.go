// Package main implements a client for Greeter service.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"
)


const (
	protocol = "unix"
	sockAddr = "/var/run/cri-resmgr/cri-resmgr-fps.sock"
)

func main() {

	str := []byte(`{"cpus": "1-3",
			"fps": "28.687267",
			"game":"stackball",
			"isFpsDrop": "false",
			"keySched": "5.90,21.79,4",
			"microTimeStamp": "1648448084957443",
			"pod": "",
			"totalSched": "6.70,49.33,238"
			}`)
	
	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)

		conn, err := net.Dial(protocol, sockAddr)
		if err != nil {
			log.Fatal(err)
		}

		
		_, err = conn.Write(str)
		if err != nil {
			log.Fatal(err)
		}

	
		err = conn.(*net.UnixConn).CloseWrite()
		if err != nil {
			log.Fatal(err)
		}

		b, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(b))
	}
}