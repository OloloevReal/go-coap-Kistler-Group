package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	coap "github.com/OloloevReal/go-coap-Kistler-Group"
)

func main() {
	co, err := coap.Dial("udp", "localhost:5688")
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	path := "/a/b"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cm, _ := co.NewGetRequest("/a")
	cm.SetQueryString("id=71747859&hash=ecd71870d1963316a97e3ac3408c9835ad8cf0f3c1bc703527c30265534f75ae")
	co.Exchange(cm)

	st := "&id=71747859&hash=ecd71870d1963316a97e3ac3408c9835ad8cf0f3c1bc703527c30265534f75ae"
	fmt.Println(st)
	strings.Split(st, "&")

	resp, err := co.GetWithContext(ctx, path)

	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	log.Printf("Response payload: %v", resp.Payload())
}
