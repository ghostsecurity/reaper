package fuzz

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gorm.io/gorm"

	"github.com/ghostsecurity/reaper/internal/database/models"
	"github.com/ghostsecurity/reaper/internal/handlers/websocket"
	"github.com/ghostsecurity/reaper/internal/types"
)

const (
	// maxWorkers        = 50
	maxWorkers        = 5
	defaultMin        = 100
	defaultMax        = 1000
	defaultMaxSuccess = 10
	defaultMaxRPS     = 10
	// TODO: input params
	// - fuzz type (header, body, param)
	// - fuzz value type (int, string, uuid)
	// - fuzz value range (int) (min, max)
	// - fuzz value format (json, urlencoded, form/string)
	// - fuzz success response codes
	// TODO: add rate limit
)

func CreateAttack(domain string, excludeKeys []string, ws *websocket.Pool, db *gorm.DB, min, max, maxSuccess int) error {
	slog.Info("Creating fuzz attack")

	// Defaults
	if min <= 0 {
		min = defaultMin
	}
	if max <= 0 {
		max = defaultMax
	}
	if maxSuccess == 0 {
		maxSuccess = defaultMaxSuccess
	}

	// Min must be less than or equal to max
	if min > max {
		min, max = max, min
	}

	req := models.Request{
		Method: http.MethodPost,
	}

	// Get most recent POST request for the domain
	res := db.Where(&req).
		Where("host LIKE ?", "%"+domain+"%").
		Order("created_at DESC").
		First(&req)

	if res.Error != nil {
		return fmt.Errorf("failed to find POST request for domain %s: %w", domain, res.Error)
	}

	slog.Info("Found request for fuzzing", "id", req.ID, "url", req.URL)

	// Parse body keys
	var bodyKeys map[string]interface{}
	if err := json.Unmarshal([]byte(req.Body), &bodyKeys); err != nil {
		return fmt.Errorf("failed to parse body keys: %w", err)
	}

	// limit concurrent workers
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxWorkers)

	var totalCount, successCount int32

	// rarly term channel
	done := make(chan struct{})
	var once sync.Once

	// iterate through body keys and fuzz each one
	for key := range bodyKeys {
		// skip excluded keys
		if slices.Contains(excludeKeys, key) {
			continue
		}
		for i := min; i <= max; i++ {
			select {
			case <-done:
				// early termination
				return nil
			default:
				wg.Add(1)
				semaphore <- struct{}{} // acquire
				go func(key string, value int) {
					defer wg.Done()
					defer func() { <-semaphore }() // release

					fuzzedReq := createFuzzedRequest(&req, key, value)
					status, err := sendRequest(fuzzedReq, ws, db)
					if err != nil {
						slog.Error("Failed to send fuzzed request", "error", err)
					} else {
						atomic.AddInt32(&totalCount, 1)
						if status >= http.StatusOK && status < http.StatusMultipleChoices {
							// inc success counter
							newCount := atomic.AddInt32(&successCount, 1)
							if int(newCount) >= maxSuccess {
								// signal early termination
								slog.Info("Max success count reached", "count", maxSuccess, "requests", totalCount)
								once.Do(func() { // Safely close the done channel
									close(done)
								})
							}
						}
					}
				}(key, i)
			}
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()

	msg := &types.AttackCompleteMessage{
		Type: types.MessageTypeAttackComplete,
	}

	ws.Broadcast <- msg

	slog.Info("Fuzz attack completed", "successCount", successCount, "totalCount", totalCount)
	return nil
}

func createFuzzedRequest(originalReq *models.Request, key string, value int) *http.Request {
	// Parse the original body
	var body map[string]interface{}
	json.Unmarshal([]byte(originalReq.Body), &body)

	// Update the specified key with the fuzzed value
	body[key] = value

	// Convert the updated body back to JSON
	fuzzedBody, _ := json.Marshal(body)

	// Create a new http.Request
	req, _ := http.NewRequest(originalReq.Method, originalReq.URL, strings.NewReader(string(fuzzedBody)))

	// Re-set headers from original request
	headerLines := strings.Split(originalReq.Headers, "\n")
	for _, line := range headerLines {
		if line = strings.TrimSpace(line); line != "" {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				headerKey := strings.TrimSpace(parts[0])
				headerValue := strings.TrimSpace(parts[1])
				req.Header.Add(headerKey, headerValue)
			}
		}
	}

	// Update Content-Length header
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(fuzzedBody)))

	return req
}

func sendRequest(req *http.Request, ws *websocket.Pool, db *gorm.DB) (status int, err error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// early return if not successful
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return resp.StatusCode, nil
	}

	requestHeaders := ""
	for key, value := range req.Header {
		requestHeaders += fmt.Sprintf("%s: %s\n", key, strings.Join(value, ","))
	}

	r, _ := req.GetBody()
	requestBody, _ := io.ReadAll(r)
	responseBody, _ := io.ReadAll(resp.Body)

	m := &types.AttackResultMessage{
		Type:      types.MessageTypeAttackResult,
		Hostname:  req.Host,
		Port:      req.URL.Port(),
		Scheme:    req.URL.Scheme,
		URL:       req.URL.String(),
		Endpoint:  req.URL.Path,
		Request:   "saved in db",
		Response:  "saved in db",
		IpAddress: req.RemoteAddr,
		Timestamp: time.Now(),
	}
	ws.Broadcast <- m

	// Create a FuzzResult and save it to the database
	fuzzResult := &models.FuzzResult{
		Hostname:  req.URL.Hostname(),
		IpAddress: req.RemoteAddr,
		Port:      req.URL.Port(),
		Scheme:    req.URL.Scheme,
		URL:       req.URL.String(),
		Endpoint:  req.URL.Path,
		Request:   string(requestHeaders) + "\n" + string(requestBody),
		Response:  string(responseBody),
	}
	res := db.Create(fuzzResult)
	if res.Error != nil {
		slog.Error("Failed to save fuzz result", "error", res.Error)
	}

	return resp.StatusCode, nil
}
