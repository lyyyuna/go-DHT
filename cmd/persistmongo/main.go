package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"godht/pkg/dht"
	"godht/pkg/mongointegr"

	"github.com/golang/glog"
)

func main() {
	port := flag.String("port", "6882", "Listening port")
	limit := flag.Int("limit", 500, "Friends made upper limit per second, default 500/seconds")
	rejoin := flag.Int("rejoin", 60, "Rejoin the DHT bootstrap rate, default 60 seconds")

	flag.Parse()
	addr := fmt.Sprintf("0.0.0.0:%s", *port)

	defer glog.Flush()

	mp, err := mongointegr.NewMongoClient("127.0.0.1:27017")
	if err != nil {
		fmt.Println(err)
		return
	}
	d, err := dht.NewDHT(addr, *limit, *rejoin)
	if err != nil {
		fmt.Println(err)
		return
	}

	d.Run()
	for {
		select {
		case a := <-d.Announcements:
			fmt.Println(a)
		case g := <-d.GetPeersQueries:
			fmt.Println(hex.EncodeToString([]byte(g.Infohash)))
			mp.InsertOneInfoHash(g)
		}
	}
}