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
 * Determine whether the request is authorized
 */
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
 * Process the webhook and log out the payload
 */
func handler(w http.ResponseWriter, r *http.Request, ntype string) {
	if authorized(w, r) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "Incorrect request")
		}

		log.WithFields(logrus.Fields{
			"url":    r.URL,
			"remote": r.RemoteAddr,
		}).Infof("Received %s\n", ntype)

		log.Info(string(body))

		w.WriteHeader(http.StatusOK)
	}
}

/*
 * Process the webhook and log out the payload
 */
func handleGeneral(w http.ResponseWriter, r *http.Request) {
	handler(w, r, "general")
}

/*
 * Process the webhook and log out the payload
 */
func handleVulnUpdate(w http.ResponseWriter, r *http.Request) {
	handler(w, r, "vuln_update")
}

/*
 * Home landing page
 */
func handleHome(w http.ResponseWriter, r *http.Request) {
	log.Info("Hit home")
	fmt.Fprintf(w, "Welcome to webhook server!")
}

func main() {
	if _, ok := os.LookupEnv("WEBHOOK_USERNAME"); !ok {
		log.Fatal("Could not find WEBHOOK_USERNAME in environment")
	}

	if _, ok := os.LookupEnv("WEBHOOK_PASSWORD"); !ok {
		log.Fatal("Could not find WEBHOOK_PASSWORD in environment")
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handleHome)
	router.HandleFunc("/v1/general", handleGeneral).Methods("POST")
	router.HandleFunc("/v1/vuln_update", handleVulnUpdate).Methods("POST")
	log.Fatal(http.ListenAndServe(":9000", router))
}
