package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type data struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type dataStore struct {
	sync.RWMutex
	data    map[string]interface{}
	path    string
	refresh time.Duration
}

func (dataStore *dataStore) GetValue(name string) interface{} {
	dataStore.RLock()
	value := dataStore.data[name]
	dataStore.RUnlock()
	return value
}

func (dataStore *dataStore) Initialize() {
	go func() {
		for {
			err := dataStore.loadData()
			if err != nil {
				log.Panic(err)
			}

			if dataStore.refresh == 0 {
				break
			}

			time.Sleep(dataStore.refresh)
		}
	}()
}

func (dataStore *dataStore) loadData() error {
	fileContents, err := ioutil.ReadFile(dataStore.path)
	if err != nil {
		return err
	}

	data := []data{}
	if json.Unmarshal(fileContents, &data); err != nil {
		return err
	}

	dataStore.Lock()
	for _, d := range data {
		dataStore.data[d.Name] = d.Value
	}
	dataStore.Unlock()

	return nil
}

func main() {
	path := flag.String("f", "data.json", "file to read JSON data from")
	port := flag.Int("p", 8081, "port the server is running on")

	dataStore := dataStore{
		data: make(map[string]interface{}),
		path: *path,
		// refresh: time.Second * 5,
	}
	dataStore.Initialize()

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			json.NewEncoder(w).Encode(dataStore.data)
		} else {
			json.NewEncoder(w).Encode(dataStore.GetValue(name))
		}
	})

	log.Printf("Data Store Server Listening on Port :%d\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), loggingHandler(router)); err != nil {
		log.Panic(err)
	}
}

func loggingHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
