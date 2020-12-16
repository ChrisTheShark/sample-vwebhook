package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// InvalidMessage will be return to the user.
	InvalidMessage = "namespace missing required team label"
	requiredLabel  = "team"
	port           = ":8443"
)

var (
	tlscert, tlskey string
)

// Namespace struct for parsing
type Namespace struct {
	Metadata Metadata `json:"metadata"`
}

// Metadata struct for parsing
type Metadata struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
}

func (m Metadata) isEmpty() bool {
	return m.Name == ""
}

// Validate handler accepts or rejects based on request contents
func Validate(w http.ResponseWriter, r *http.Request) {
	arReview := v1beta1.AdmissionReview{}
	if err := json.NewDecoder(r.Body).Decode(&arReview); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if arReview.Request == nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	raw := arReview.Request.Object.Raw

	ns := Namespace{}
	if err := json.Unmarshal(raw, &ns); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if ns.Metadata.isEmpty() {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	arReview.Response = &v1beta1.AdmissionResponse{
		UID:     arReview.Request.UID,
		Allowed: true,
	}

	if len(ns.Metadata.Labels) == 0 || ns.Metadata.Labels[requiredLabel] == "" {
		arReview.Response.Allowed = false
		arReview.Response.Result = &metav1.Status{
			Message: InvalidMessage,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&arReview)
}

func main() {
	flag.StringVar(&tlscert, "tlsCertFile", "/etc/certs/tls.crt",
		"File containing a certificate for HTTPS.")
	flag.StringVar(&tlskey, "tlsKeyFile", "/etc/certs/tls.key",
		"File containing a private key for HTTPS.")
	flag.Parse()

	http.HandleFunc("/validate", Validate)
	log.Fatal(http.ListenAndServeTLS(port, tlscert, tlskey, nil))
}
