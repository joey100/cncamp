package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"net/http/pprof"
	_ "net/http/pprof"

	"github.com/cnncamp-homework/joey/homework/cncamp/httpserver/metrics"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("Starting http server...")
	//http.HandleFunc("/", rootHandler)
	c, python, java := true, false, "no!"
	fmt.Println(c, python, java)
	//err := http.ListenAndServe(":80", nil)

	metrics.Register()
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/metrics", promhttp.Handler())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    ":80",
		Handler: mux,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Print("Server started.")

	<-done
	log.Print("Server Stopped.")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")

}

func healthz(w http.ResponseWriter, r *http.Request) {
	//io.WriteString(w, "ok\n")
	//w.Header().Set("Location", "https://www.baidu.com")
	//w.WriteHeader(301)
	// v := os.Getenv("VERSION")
	// if v == "" {
	// 	v = "VERSION NOT FOUND"
	// }
	// w.Header().Add("VERSION", v)
	// for k, v := range r.Header {
	// 	stringV := strings.Join(v, ",")
	// 	w.Header().Add(k, stringV)
	// }
	w.Write([]byte("OK"))
	fmt.Println(time.Now(), r.URL, http.StatusOK)
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	delay := randInt(10, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))
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
