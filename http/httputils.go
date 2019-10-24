package httputils

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httputil"

	stringutils "github.com/alessiosavi/GoGPUtils/string"
)

// CreateCookie is delegated to initialize and set a cookie
func CreateCookie(name, value, domain, path string, maxage int, httpOnly bool) (*http.Cookie, error) {

	if stringutils.IsBlank(name) {
		return nil, errors.New("cookie name not provided")
	}
	if stringutils.IsBlank(value) {
		return nil, errors.New("cookie value not provided")
	}
	// Override path in case of not provided
	if path == "" {
		path = "/"
	}

	// Set session cookie in case of lesser than 0
	if maxage < 0 {
		maxage = 0
	}

	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Domain:   domain,
		Path:     path,
		MaxAge:   maxage,
		HttpOnly: httpOnly,
	}

	return &cookie, nil
}

// ServeCookie is delegated to set a cookie and shutdown the server after the first call
func ServeCookie(ip, port, endpoint, name, value, domain, path string, maxage int, httpOnly bool) error {

	// Validate input
	if stringutils.IsBlank(ip) {
		return errors.New("hostname/ip not provided")
	}
	if stringutils.IsBlank(port) {
		return errors.New("port not provided")
	}
	if stringutils.IsBlank(endpoint) {
		endpoint = "/"
	}

	// Instantiate server
	m := http.NewServeMux()
	s := http.Server{Addr: ip + ":" + port, Handler: m}

	// Create cookie
	cookie, err := CreateCookie(name, value, domain, path, maxage, httpOnly)
	if err != nil {
		return err
	}

	// Bind the endpoint for instantiate the cookie
	m.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, cookie)
		w.WriteHeader(200)
		go http.Get(`http://` + ip + `:` + port + `/shutdown`)
	})

	// Bind the endpoint for shutdown the server
	m.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		s.Shutdown(context.Background())
	})

	// Serve the http webserver
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
		return err
	}

	return nil
}

// DebugRequest is delegated to print the request data for debug purpouse
func DebugRequest(ip, port, endpoint string) error {

	// Validate input
	if stringutils.IsBlank(ip) {
		return errors.New("hostname/ip not provided")
	}
	if stringutils.IsBlank(port) {
		return errors.New("port not provided")
	}
	if stringutils.IsBlank(endpoint) {
		endpoint = "/"
	}

	// Instantiate server
	m := http.NewServeMux()
	s := http.Server{Addr: ip + ":" + port, Handler: m}
	// Bind the endpoint for instantiate the cookie
	m.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		data, err := httputil.DumpRequest(r, true)
		log.Println("Request data -> \n", string(data))
		if err != nil {
			log.Println("Errors -> ", err)
		}
		w.Write(data)
		go http.Get(`http://` + ip + `:` + port + `/shutdown`)
	})

	// Bind the endpoint for shutdown the server
	m.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		s.Shutdown(context.Background())
	})

	// Serve the http webserver
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
		return err
	}

	return nil
}
