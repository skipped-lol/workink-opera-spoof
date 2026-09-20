package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func main() {
	fmt.Println("Spoofing Opera task...")

	client := &http.Client{}

	println("Sending request #1")
	req, err := http.NewRequest("GET", "https://work.ink/_api/v2/affiliate/operaGX", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("User-Agent", "Opera Installer/1.0")
	req.Header.Set("Host", "work.ink")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	//fmt.Printf("Response body: %s\n", body)

	println("Sending request #2")
	jsonData := []byte(`{"noteligible":true}`)

	req, err = http.NewRequest(
		"POST",
		"https://work.ink/_api/v2/callback/operaGX",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Opera Installer/1.0")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Printf("Response body: %s\n", body)
}
