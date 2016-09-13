package tigerblood

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// IPAddressFromHTTPPath takes a HTTP path and returns an IPv4 IP if it's found, or an error if none is found.
func IPAddressFromHTTPPath(path string) (string, error) {
	path = path[1:len(path)]
	ip, network, err := net.ParseCIDR(path)
	if err != nil {
		if strings.Contains(path, "/") {
			return "", fmt.Errorf("Error getting IP from HTTP path: %s", err)
		}
		ip = net.ParseIP(path)
		if ip == nil {
			return "", fmt.Errorf("Error getting IP from HTTP path: %s", err)
		}
		network = &net.IPNet{}
		if ip.To4() != nil {
			network.Mask = net.CIDRMask(32, 32)
		} else if ip.To16() != nil {
			network.Mask = net.CIDRMask(128, 128)
		}
	}
	network.IP = ip
	return network.String(), nil
}

// Handler is the main HTTP handler for tigerblood.
func Handler(w http.ResponseWriter, r *http.Request, db *DB) {
	startTime := time.Now()
	switch r.Method {
	case "GET":
		ReadReputation(w, r, db)
	case "POST":
	case "PUT":
	case "DELETE":
	default:
	}
	if time.Since(startTime).Nanoseconds() > 1e7 {
		log.Printf("Request took %s to proces\n", time.Since(startTime))
	}
}

// CreateReputation takes a JSON formatted IP reputation entry from the http request and inserts it to the database.
func CreateReputation(w http.ResponseWriter, r *http.Request) {
}

// ReadReputation returns a JSON-formatted reputation entry from the database.
func ReadReputation(w http.ResponseWriter, r *http.Request, db *DB) {
	ip, err := IPAddressFromHTTPPath(r.URL.Path)
	if err != nil {
		// This means there was no IP address found in the path
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("No IP address found in path %s: %s", r.URL.Path, err)
		return
	}
	entry, err := db.SelectSmallestMatchingSubnet(ip)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("No entries found for IP %s", ip)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error executing SQL: %s", err)
	}
	json, err := json.Marshal(entry)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error marshaling JSON: %s", err)
	}
	w.Write(json)
}

// UpdateReputation takes a JSON body from the http request and updates that reputation entry on the database.
func UpdateReputation(w http.ResponseWriter, r *http.Request) {

}

// DeleteReputation deletes
func DeleteReputation(w http.ResponseWriter, r *http.Request) {

}
