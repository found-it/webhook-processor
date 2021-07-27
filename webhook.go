package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var logging = logrus.New()
var log = logging.WithFields(logrus.Fields{"server": "0.0.0.0:9000"})

/*
 *  [ Handler ] Home landing page
 */
func homeLink(w http.ResponseWriter, r *http.Request) {
	log.Info("Hit home")
	fmt.Fprintf(w, "Welcome to webhook server!")
}

func authorized(w http.ResponseWriter, r *http.Request) bool {
	if u, p, ok := r.BasicAuth(); ok {
		if u == os.Getenv("WEBHOOK_USERNAME") && p == os.Getenv("WEBHOOK_PASSWORD") {
			return true
		}
		log.WithFields(logrus.Fields{
			"username": u,
		}).Error("Unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorized")
		return false
	}

	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, "Unauthorized")
	log.Error("Parsing basic auth failed")
	return false
}

/*

 */
func processor(w http.ResponseWriter, r *http.Request) {

	if authorized(w, r) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "Incorrect request")
		}

		log.WithFields(logrus.Fields{
			"url": r.URL,
		}).Info("Received!")

		log.Info(string(body))

		w.WriteHeader(http.StatusOK)
	}
}

/*
 *  Use Gorilla Mux to handle routes
 */
func main() {
	if _, ok := os.LookupEnv("WEBHOOK_USERNAME"); !ok {
		log.Fatal("Could not find WEBHOOK_USERNAME in environment")
	}

	if _, ok := os.LookupEnv("WEBHOOK_PASSWORD"); !ok {
		log.Fatal("Could not find WEBHOOK_PASSWORD in environment")
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/v1/webhook", processor).Methods("POST")
	log.Fatal(http.ListenAndServe(":9000", router))
}
