package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: 2 * time.Second,
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: go run ./klien-go <address> <see|add> [argumen...]\n")
		os.Exit(1)
	}

	baseURL := strings.TrimSuffix(args[0], "/")

	if len(args) == 1 {
		cmdAllList(baseURL)
		return
	}

	command := args[1]

	switch command {
	case "lihat":
		if len(args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: go run ./klien-go <address> see <id>\n")
			os.Exit(1)
		}
		id := args[2]
		cmdSeeList(baseURL, id)

	case "tambah":
		if len(args) < 4 {
			fmt.Fprintf(os.Stderr, "Usage: go run ./klien-go <address> add <nama> <lokasi>\n")
			os.Exit(1)
		}
		nama := args[2]
		lokasi := args[3]
		cmdMakeNewList(baseURL, nama, lokasi)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func cmdAllList(baseURL string) {
	url := fmt.Sprintf("%s/barang", baseURL)
	maxRetries := 5
	backoff := 200 * time.Millisecond

	for i := 1; i <= maxRetries; i++ {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating request: %v\n", err)
			os.Exit(1)
		}

		resp, err := client.Do(req)
		if err != nil {
			if i == maxRetries {
				fmt.Fprintf(os.Stderr, "Connection error after %d tries: %v\n", maxRetries, err)
				os.Exit(1)
			}
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			if i == maxRetries {
				fmt.Fprintf(os.Stderr, "Failed to read response body: %v\n", readErr)
				os.Exit(1)
			}
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			if i == maxRetries {
				fmt.Println(string(body))
				os.Exit(1)
			}
			retryAfterSec := 0
			if headerVal := resp.Header.Get("Retry-After"); headerVal != "" {
				if sec, parseErr := strconv.Atoi(headerVal); parseErr == nil {
					retryAfterSec = sec
				}
			}

			if retryAfterSec > 0 {
				time.Sleep(time.Duration(retryAfterSec) * time.Second)
			} else {
				time.Sleep(backoff)
				backoff *= 2
			}
			continue
		}

		fmt.Println(string(strings.TrimSpace(string(body))))
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

func cmdSeeList(baseURL, id string) {
	url := fmt.Sprintf("%s/barang/%s", baseURL, id)
	maxRetries := 5
	backoff := 200 * time.Millisecond

	for i := 1; i <= maxRetries; i++ {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating request: %v\n", err)
			os.Exit(1)
		}

		resp, err := client.Do(req)
		if err != nil {
			if i == maxRetries {
				fmt.Fprintf(os.Stderr, "Connection error after %d tries: %v\n", maxRetries, err)
				os.Exit(1)
			}
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			if i == maxRetries {
				fmt.Fprintf(os.Stderr, "Failed to read response body: %v\n", readErr)
				os.Exit(1)
			}
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			if i == maxRetries {
				fmt.Println(string(body))
				os.Exit(1)
			}
			retryAfterSec := 0
			if headerVal := resp.Header.Get("Retry-After"); headerVal != "" {
				if sec, parseErr := strconv.Atoi(headerVal); parseErr == nil {
					retryAfterSec = sec
				}
			}

			if retryAfterSec > 0 {
				time.Sleep(time.Duration(retryAfterSec) * time.Second)
			} else {
				time.Sleep(backoff)
				backoff *= 2
			}
			continue
		}
		fmt.Println(string(strings.TrimSpace(string(body))))
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

func cmdMakeNewList(baseURL, nama, lokasi string) {
	url := fmt.Sprintf("%s/barang", baseURL)
	payload := map[string]string{
		"nama":   nama,
		"lokasi": lokasi,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal JSON: %v\n", err)
		os.Exit(1)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application /json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to do Request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read response: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(strings.TrimSpace(string(body))))
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
