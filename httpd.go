package main

import (
	"fmt"
	"net/http"
	"syscall"
)

const portMaxNum = 65535

// Httpd struct expressses http server that listen in Port and trigger Callback
type Httpd struct {
	Port     uint
	Prefix   string
	Callback func(sig syscall.Signal) error
}

// NewHttpd creates new httpd instance with option
func NewHttpd(port uint, prefix string, callback func(sig syscall.Signal) error) (*Httpd, error) {
	if port > portMaxNum {
		return nil, fmt.Errorf("port number shold be smaller than %d", portMaxNum)
	}

	return &Httpd{
		Port:     port,
		Prefix:   prefix,
		Callback: callback,
	}, nil
}

// Run function starts http server
// http server handle signal-path generated by metadata of signals in signal.go
func (httpd *Httpd) Run() {
	for i := 0; i < len(Signals); i++ {
		signal := Signals[i]
		signalStr := SignalStrs[i]
		http.HandleFunc(httpd.Prefix+"/"+signalStr, func(w http.ResponseWriter, r *http.Request) {
			err := httpd.Callback(signal)
			if err != nil {
				// If failed to process Callback, http server returns 500
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Failed to poroxy %s to destination command", signalStr)
			} else {
				fmt.Fprintf(w, "Successed to proxy %s to destination command", signalStr)
			}
		})
	}

	http.ListenAndServe(fmt.Sprintf(":%d", httpd.Port), nil)
}
