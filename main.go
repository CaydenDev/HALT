package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Challenge struct {
	ID        string
	A         *big.Int
	B         *big.Int
	Result    *big.Int
	Timestamp time.Time
}

type ProxyServer struct {
	target     *url.URL
	proxy      *httputil.ReverseProxy
	challenges sync.Map
}

func NewProxyServer(targetURL string) (*ProxyServer, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	return &ProxyServer{
		target: target,
		proxy:  httputil.NewSingleHostReverseProxy(target),
	}, nil
}

func (p *ProxyServer) generateChallenge() *Challenge {
	a, _ := rand.Int(rand.Reader, big.NewInt(100))
	b, _ := rand.Int(rand.Reader, big.NewInt(100))
	result := new(big.Int).Add(a, b)

	id := make([]byte, 16)
	rand.Read(id)

	challenge := &Challenge{
		ID:        base64.URLEncoding.EncodeToString(id),
		A:         a,
		B:         b,
		Result:    result,
		Timestamp: time.Now(),
	}

	p.challenges.Store(challenge.ID, challenge)
	return challenge
}

func (p *ProxyServer) verifyChallenge(id string, answer string) bool {
	val, exists := p.challenges.Load(id)
	if !exists {
		return false
	}

	challenge := val.(*Challenge)
	p.challenges.Delete(id)

	if time.Since(challenge.Timestamp) > 5*time.Minute {
		return false
	}

	userAnswer := new(big.Int)
	userAnswer.SetString(answer, 10)
	return userAnswer.Cmp(challenge.Result) == 0
}

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("halt-verified"); err == nil && cookie.Value == "true" {
		p.proxy.ServeHTTP(w, r)
		return
	}

	if r.Method == "POST" && r.URL.Path == "/verify" {
		id := r.FormValue("challenge_id")
		answer := r.FormValue("answer")

		if p.verifyChallenge(id, answer) {
			http.SetCookie(w, &http.Cookie{
				Name:     "halt-verified",
				Value:    "true",
				Path:     "/",
				MaxAge:   3600,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/challenge", http.StatusSeeOther)
		return
	}

	challenge := p.generateChallenge()
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>HALT - Human Authentication</title>
			<style>
				body { font-family: Arial, sans-serif; margin: 40px; }
				.challenge { background: #f5f5f5; padding: 20px; border-radius: 5px; }
			</style>
		</head>
		<body>
			<div class="challenge">
				<h2>Please solve this challenge to continue</h2>
				<p>What is %s + %s?</p>
				<form method="POST" action="/verify">
					<input type="hidden" name="challenge_id" value="%s">
					<input type="number" name="answer" required>
					<button type="submit">Submit</button>
				</form>
			</div>
		</body>
		</html>
	`, challenge.A, challenge.B, challenge.ID)
}

func main() {
	targetURL := "http://localhost:8080"
	proxyPort := 3000

	proxy, err := NewProxyServer(targetURL)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("HALT proxy server starting on port %d, forwarding to %s", proxyPort, targetURL)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", proxyPort), proxy))
}