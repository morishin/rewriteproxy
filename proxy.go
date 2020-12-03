package main

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/gobwas/glob"
)

func ReverseProxy(webAppUrl *url.URL, cloudFunctionBaseUrl *url.URL, firebaseJson FirebaseJson) *httputil.ReverseProxy {
	isPathToRedirect := func(path string) bool {
		return path == "/index.html" || firebaseJson.Hosting.CleanUrls && strings.Contains(path, ".html")
	}

	// Setip redirect server
	redirectServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/index.html" {
			// Redirect `/index.html` to `/``
			res.Header().Set("Location", "/")
			res.WriteHeader(http.StatusMovedPermanently)
		} else if firebaseJson.Hosting.CleanUrls && strings.Contains(req.URL.Path, ".html") {
			// Redirect `/some.html` to `/some`
			cleanUrl := strings.Replace(req.URL.Path, ".html", "", 1)
			res.Header().Set("Location", cleanUrl)
			res.WriteHeader(http.StatusMovedPermanently)
		}
	}))
	redirectServerUrl, _ := url.Parse(redirectServer.URL)

	// Setup glob and regexp representing Firebase Hosting rewrite rules
	globs := map[string]glob.Glob{}
	regexps := map[string]*regexp.Regexp{}
	for _, rule := range firebaseJson.Hosting.Rewrites {
		if rule.Source != "" {
			globs[rule.Source] = glob.MustCompile(rule.Source, '.')
		} else if rule.Regex != "" {
			regexps[rule.Regex] = regexp.MustCompile(rule.Regex)
		}
	}

	// Setup director which rewrites host and path of requests
	director := func(req *http.Request) {
		if isPathToRedirect(req.URL.Path) {
			// Pass requests to redirect server
			req.URL.Scheme = redirectServerUrl.Scheme
			req.URL.Host = redirectServerUrl.Host
			return
		}

		for _, rule := range firebaseJson.Hosting.Rewrites {
			if rule.Source != "" && globs[rule.Source].Match(req.URL.Path) || rule.Regex != "" && regexps[rule.Regex].MatchString(req.URL.Path) {
				if rule.Destination != "" {
					// Proxy to web app
					req.URL.Scheme = webAppUrl.Scheme
					req.URL.Host = webAppUrl.Host
					req.URL.Path = rule.Destination
				} else if rule.Function != "" {
					// Proxy to Cloud Function HTTP trigger
					req.URL.Scheme = cloudFunctionBaseUrl.Scheme
					req.URL.Host = cloudFunctionBaseUrl.Host
					req.URL.Path = singleJoiningSlash(cloudFunctionBaseUrl.Path, rule.Function)
				}
				return
			}
		}

		// Proxy to web app
		req.URL.Scheme = webAppUrl.Scheme
		req.URL.Host = webAppUrl.Host
		if req.URL.Path == "/" {
			req.URL.Path = "/index.html"
		} else if !strings.Contains(req.URL.Path, ".") {
			req.URL.Path += ".html"
		}
	}

	return &httputil.ReverseProxy{Director: director}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
