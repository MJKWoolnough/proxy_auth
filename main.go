// +build ignore

package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	auth "github.com/MJKWoolnough/proxy_auth"
)

var (
	port    = flag.String("p", "80", "server port")
	address = flag.String("a", "", "server domain")
	data    = flag.String("u", "", "site/user json file")
)

func main() {
	flag.Parse()
	f, err := os.Open(*data)
	if err != nil {
		log.Println(err)
		return
	}
	err = auth.Setup(f, http.DefaultServeMux)
	f.Close()
	if err != nil {
		log.Println(err)
		switch errs := err.(type) {
		case auth.ProcessErrors:
			for _, e := range errs {
				log.Println(e)
			}
		}
		return
	}
	log.Println("Server starting: " + *address + ":" + *port)
	log.Println("API Path: " + "http://" + *address + ":" + *port + auth.APIPath + "{domain}" + auth.APIFile)
	log.Fatal(http.ListenAndServe(*address+":"+*port, nil))
}
