package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"

	coap "github.com/OloloevReal/go-coap-Kistler-Group"
)

func handleA(w coap.ResponseWriter, req *coap.Request) {
	log.Printf("Got message in handleA: path=%q: %#v from %v", req.Msg.Path(), req.Msg, req.Client.RemoteAddr())
	w.SetContentFormat(coap.TextPlain)
	log.Printf("Transmitting from A")
	ctx, cancel := context.WithTimeout(req.Ctx, time.Second)
	defer cancel()
	if _, err := w.WriteWithContext(ctx, []byte("hello world")); err != nil {
		log.Printf("Cannot send response: %v", err)
	}
}

func handleB(w coap.ResponseWriter, req *coap.Request) {
	log.Printf("Got message in handleB: path=%q: %#v from %v", req.Msg.Path(), req.Msg, req.Client.RemoteAddr())
	resp := w.NewResponse(coap.Content)
	resp.SetOption(coap.ContentFormat, coap.TextPlain)
	resp.SetPayload([]byte("good bye!"))
	log.Printf("Transmitting from B %#v", resp)
	ctx, cancel := context.WithTimeout(req.Ctx, time.Second)
	defer cancel()
	if err := w.WriteMsgWithContext(ctx, resp); err != nil {
		log.Printf("Cannot send response: %v", err)
	}
}

func handleAB(w coap.ResponseWriter, req *coap.Request) {
	//log.Printf("Got message in handleAB: path=%q: %#v from %v", req.Msg.Path(), req.Msg, req.Client.RemoteAddr())
	//w.SetContentFormat(coap.TextPlain)
	log.Printf("Transmitting from AB")
	ctx, cancel := context.WithTimeout(req.Ctx, time.Second)
	defer cancel()
	if _, err := w.WriteWithContext(ctx, nil); err != nil {
		log.Printf("Cannot send response: %v", err)
	}

	// if _, err := w.WriteWithContext(ctx, []byte("hello world")); err != nil {
	// 	log.Printf("Cannot send response: %v", err)
	// }
}

func main() {
	mux := coap.NewServeMux()
	mux.Use(Recover)
	mux.Use(Logs)
	mux.Use(Auth)
	mux.Handle("/a", Logs(coap.HandlerFunc(handleA)))
	mux.Handle("/a", coap.HandlerFunc(handleA))
	mux.Handle("/b", coap.HandlerFunc(handleB))
	mux.Handle("/data", coap.HandlerFunc(handleAB))
	mux.Handle("/data/get", coap.HandlerFunc(handleAB))

	log.Fatal(coap.ListenAndServe("udp", ":5683", mux))
}

func Recover(next coap.Handler) coap.Handler {
	fn := func(w coap.ResponseWriter, r *coap.Request) {
		defer func() {
			//log.Println("defer recover")

			if rec := recover(); rec != nil {
				fmt.Fprintf(os.Stderr, "Panic: %+v\n", rec)
				debug.PrintStack()
			}
			newMesg := w.NewResponse(coap.InternalServerError)
			r.Client.Exchange(newMesg)
		}()
		next.ServeCOAP(w, r)
	}
	return coap.HandlerFunc(fn)
}

func Logs(next coap.Handler) coap.Handler {
	fn := func(w coap.ResponseWriter, req *coap.Request) {
		log.Printf("\"%s://%s/%s\" from %s\t%#v\r\n", "coap", req.Client.LocalAddr().String(), req.Msg.PathString(), req.Client.RemoteAddr().String(), req.Msg)
		next.ServeCOAP(w, req)
	}
	return coap.HandlerFunc(fn)
}

func Auth(next coap.Handler) coap.Handler {
	fn := func(w coap.ResponseWriter, req *coap.Request) {
		hash, id := "", ""

		aq := req.Msg.Query()
		for _, q := range aq {
			qs := strings.Split(q, "=")
			if len(qs) == 2 {
				switch qs[0] {
				case "id":
					id = qs[1]
				case "hash":
					hash = qs[1]
				}
			}
		}
		if !checkAuth(id, hash) {
			log.Println("Auth failed")
			newMesg := w.NewResponse(coap.Unauthorized)
			req.Client.Exchange(newMesg)
			return
		}
		next.ServeCOAP(w, req)
	}
	return coap.HandlerFunc(fn)
}

func checkAuth(id string, hash string) bool {
	if id == "71747859" && hash == "ecd71870d1963316a97e3ac3408c9835ad8cf0f3c1bc703527c30265534f75ae" {
		return true
	}
	return false
}
