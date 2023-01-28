package main

import (
	"crypto/tls"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"punyshort-redirect-proxy/apiclient"
	"punyshort-redirect-proxy/helper"
	"strings"
	"time"
)

func main() {
	println("Starting...")
	apiClient := apiclient.NewClient(os.Getenv("PUNYSHORT_BASE_URL"), os.Getenv("PUNYSHORT_KEY"))
	errorUrl := os.Getenv("PUNYSHORT_ERROR_URL")

	followRedirect := func(writer http.ResponseWriter, request *http.Request) {
		ip, _ := helper.GetIP(request, os.Getenv("PUNYSHORT_IP_FORWARDING") == "true")

		shorten, err := apiClient.FollowRedirection(apiclient.RedirectionData{
			Domain:    request.Host,
			Path:      request.URL.Path[1:],
			Ip:        ip,
			UserAgent: request.UserAgent(),
			Referrer:  request.Referer(),
		})

		if err != nil {
			log.Println(err)

			writer.Header().Set("Location", errorUrl+"?error=Internal")
			writer.WriteHeader(307)
			return
		}

		if shorten.Error {
			if errorUrl != "" {
				writer.Header().Set("Location", errorUrl+"?error="+shorten.Exception)
				writer.WriteHeader(307)
			} else {
				writer.Write([]byte(errorUrl + "Error: " + strings.Replace(shorten.Exception, "Exception", "", 1)))
			}
			return
		}

		writer.Header().Set("Location", shorten.LongLink)
		writer.WriteHeader(307)
	}

	useSSL := false

	if useSSL && os.Getenv("PUNYSHORT_USE_SSL") == "true" {
		useSSL = true
	}

	mux := http.NewServeMux()

	domains := []string{}

	mux.HandleFunc("/", followRedirect)

	certManager := &autocert.Manager{}

	if useSSL {
		certManager.Prompt = autocert.AcceptTOS
		certManager.Cache = autocert.DirCache("letsencrypt")

		server := http.Server{
			Addr:    ":443",
			Handler: mux,
			TLSConfig: &tls.Config{GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {

				println("Serving " + info.ServerName)
				for _, name := range domains {
					if name == info.ServerName {
						return certManager.GetCertificate(info)
					}
				}

				println("Adding " + info.ServerName)
				domains = append(domains, info.ServerName)
				certManager.HostPolicy = autocert.HostWhitelist(domains...)

				return certManager.GetCertificate(info)
			}},
			ReadHeaderTimeout: 60 * time.Second,
		}

		go server.ListenAndServeTLS("", "")
	}

	httpMux := http.NewServeMux()

	httpMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		followRedirect(writer, request)

		if useSSL {
			host := request.Host

			for _, name := range domains {
				if name == host {
					return
				}
			}
			domains = append(domains, host)
			certManager.HostPolicy = autocert.HostWhitelist(domains...)
		}
	})

	log.Fatal(http.ListenAndServe(":80", httpMux))
}
