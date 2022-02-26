package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"net/http/pprof"
	_ "net/http/pprof"

	"github.com/golang/glog"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("Starting http server...")
	//http.HandleFunc("/", rootHandler)
	c, python, java := true, false, "no!"
	fmt.Println(c, python, java)
	//err := http.ListenAndServe(":80", nil)
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatal(err)
	}

}

func healthz(w http.ResponseWriter, r *http.Request) {
	//io.WriteString(w, "ok\n")
	//w.Header().Set("Location", "https://www.baidu.com")
	//w.WriteHeader(301)
	v := os.Getenv("VERSION")
	if v == "" {
		v = "VERSION NOT FOUND"
	}
	w.Header().Add("VERSION", v)
	for k, v := range r.Header {
		stringV := strings.Join(v, ",")
		w.Header().Add(k, stringV)
	}
	w.Write([]byte("OK"))
	fmt.Println(time.Now(), r.URL, http.StatusOK)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	ip, _ := GetIP(r)
	v := os.Getenv("VERSION")
	if v == "" {
		v = "VERSION NOT FOUND"
	}
	w.Header().Add("VERSION", v)
	for k, v := range r.Header {
		//io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
		stringV := strings.Join(v, ",")
		w.Header().Add(k, stringV)
	}
	headerString, _ := json.Marshal(r.Header)
	statusCode := 200
	if r.URL.Path != "/" {
		statusCode = 404
		w.Write([]byte("URL not found. Please double check."))
	} else {
		w.Write([]byte("Welcome to CNCAMP's http server."))
	}
	fmt.Println(time.Now(), ip, r.URL, statusCode, string(headerString))

}

// GetIP returns request real ip.
func GetIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	return "", errors.New("no valid ip found")
}
