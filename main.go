package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"golang.org/x/crypto/bcrypt"
)

const (
	version = "v2.2.1"
)

func genLoad2(w http.ResponseWriter) {
	list := WordList()
	sort.Strings(list)
	wordCount := make(map[string]int)
	log.Println("Generating Load")
	for _, item := range list {
		wordCount[item] = len(item)
	}
	log.Println("Load Generation Stopped.")
	w.WriteHeader(200)
}
func genLoad(w http.ResponseWriter, r *http.Request) {
	password := []byte("asd907234han!2hjfads.")
	var wg sync.WaitGroup
	log.Println(r.RequestURI)
	URIParts := strings.Split(r.RequestURI, "/")
	var count int
	var err error
	if len(URIParts) == 3 {
		count, err = strconv.Atoi(URIParts[2])
		if err != nil {
			log.Println(err)
		}
	} else {
		count = 2
	}

	wg.Add(count)
	// Hashing the password with the default cost of 10
	log.Println("Generating Load")
	for i := 0; i < count; i++ {
		go func() {
			_, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	log.Println("Load Generation Stopped.")
	w.WriteHeader(200)

}

func main() {
	log.Println("Starting application...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world!\n")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, version)
	})

	http.HandleFunc("/updates", func(w http.ResponseWriter, r *http.Request) {
		contents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("%s", err)
		}
		log.Printf("%s\n", string(contents))
	})

	http.HandleFunc("/load/", func(w http.ResponseWriter, r *http.Request) {
		genLoad2(w)
	})

	s := http.Server{Addr: ":8080"}
	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")

	s.Shutdown(context.Background())
}
