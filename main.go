package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "v1.0.0"

func main() {
	// log with time and line
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	log.Println(version)
	log.Println("this is biz server")

	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", Hello)
	log.Println("register route : " + "/echo")
	mux.HandleFunc("/readyz", Ping)
	log.Println("register route : " + "/readyz")
	mux.HandleFunc("/healthz", Ping)
	log.Println("register route : " + "/healthz")

	mux.HandleFunc("/", Default)
	log.Println("register route : " + "/")

	log.Print("start serving ... ")
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		panic(err)
	}

	log.Print("end serving ... ")
}

func Default(w http.ResponseWriter, r *http.Request) {
	// return namespace if in k8s
	// if in k8s
	if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		// read namespace
		namespace, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
		if err != nil {
			log.Println("read namespace err : " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		content := "this is biz server in k8s, namespace : " + string(namespace)
		// pod name
		podName, err := os.Hostname()
		if err != nil {
			log.Println("get pod name err : " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		content += ", pod name : " + podName
		content += ", version: " + version

		w.Write([]byte(content))
		return
	} else {
		w.Write([]byte("this is biz server, version: " + version))
		return
	}

}

func checkSidecar() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		"http://localhost:9000/echo",
		bytes.NewReader([]byte("this is biz")))
	if err != nil {
		return
	}
	defer request.Body.Close()
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		err = errors.New("sidecar not ready")
		return
	}

	body, _ := io.ReadAll(response.Body)

	log.Println("response from sidecar : " + string(body))
	return
}

// Ping returns true automatically when checked.
var Ping = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

var Hello = func(w http.ResponseWriter, r *http.Request) {
	// write response
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("hello world"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
