package main // import "github.com/zgiles/yesserver/cmd/yesserver"

import (
	"flag"
	"fmt"
	"log"
	"os"

	"layeh.com/radius"
	. "layeh.com/radius/rfc2865"
)

var secret = flag.String("secret", "", "shared RADIUS secret between clients and server")

func handler(w radius.ResponseWriter, r *radius.Request) {
	username, err1 := UserName_LookupString(r.Packet)
	_, err2 := UserPassword_LookupString(r.Packet)
	if err1 != nil || err2 != nil {
		w.Write(r.Response(radius.CodeAccessReject))
		return
	}
	log.Printf("%s requesting access (%s #%d)\n", username, r.RemoteAddr, r.Identifier)
	var code radius.Code
	output := []byte("HOPE")
	code = radius.CodeAccessAccept
	log.Printf("%s accepted (%s #%d)\n", username, r.RemoteAddr, r.Identifier)
	resp := r.Response(code)
	if len(output) > 0 {
		ReplyMessage_Set(r.Packet, output)
	}
	w.Write(resp)
}

const usage = `
Just accepts any packet with a username and password in it
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [flags]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, usage)
	}
	flag.Parse()

	if *secret == "" {
		flag.Usage()
		os.Exit(1)
	}

	log.Println("radserver starting")

	server := radius.PacketServer{
		Handler:      radius.HandlerFunc(handler),
		SecretSource: radius.StaticSecretSource([]byte(*secret)),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
