package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
)

func main() {
	fmt.Println("Spoofing Opera task...")

	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("Cookie error:", err)
		return
	}

	transport := &http.Transport{
		ForceAttemptHTTP2: false,
		TLSNextProto:      make(map[string]func(string, *tls.Conn) http.RoundTripper),
	}

	client := &http.Client{
		Jar:       jar,
		Transport: transport,
	}

	fmt.Println("Sending request #1")

	req, err := http.NewRequest(
		"GET",
		"https://work.ink/_api/v2/affiliate/operaGX",
		nil,
	)
	if err != nil {
		fmt.Println("Error on request:", err)
		return
	}

	req.Header.Set("User-Agent", "Opera Installer/1.0")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error on request:", err)
		return
	}
	resp.Body.Close()

	fmt.Println("Response #1:", resp.Status)

	fmt.Println("\nSending request #2")

	jsonData := []byte(`{"noteligible":true}`)

	req, err = http.NewRequest(
		"POST",
		"https://work.ink/_api/v2/callback/operaGX",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		fmt.Println("Error on request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Opera Installer/1.0")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error on request:", err)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Println("Response #2:", resp.Status)
	fmt.Println("Response #2 body:", string(body))
}
