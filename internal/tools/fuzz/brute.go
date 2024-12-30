package fuzz

import (
    "encoding/json"
    "fmt"
    "io"
    "log/slog"
    "net/http"
    "strings"
    "sync"
    "sync/atomic"
    "time"

    "gorm.io/gorm"

    "github.com/ghostsecurity/reaper/internal/database/models"
    "github.com/ghostsecurity/reaper/internal/handlers/websocket"
    "github.com/ghostsecurity/reaper/internal/types"
)

// CreateBruteForceAttack attempts all values in 'dictionary' for a single key/param
func CreateBruteForceAttack(
    attackID uint,
    hostname string,
    param string,
    dictionary []string,       // e.g. {"000","001","002",...,"999"}
    ws *websocket.Pool,
    db *gorm.DB,
    maxSuccess int,            // e.g. 10
    maxRPS int,                // e.g. 10
) error {
    slog.Info("Creating brute force attack")

    // Default fallbacks if not provided
    if maxSuccess <= 0 {
        maxSuccess = 10
    }
    if maxRPS <= 0 {
        maxRPS = 10
    }

    // Find a request to base our brute force attempts on
    baseReq := models.Request{
        Method: http.MethodPost, // or whatever method
    }
    res := db.Where(&baseReq).
        Where("host LIKE ?", "%"+hostname+"%").
        Order("created_at DESC").
        First(&baseReq)

    if res.Error != nil {
        return fmt.Errorf("failed to find request for hostname %s: %w", hostname, res.Error)
    }

    slog.Info("Found request for brute force", "id", baseReq.ID, "url", baseReq.URL)

    // We'll need concurrency control
    var wg sync.WaitGroup
    var activeWorkers int32
    var successCount int32
    var totalCount int32

    done := make(chan struct{})
    var once sync.Once

    // Convert requests per second into a simple concurrency limit
    // so we don't hammer the server too fast
    semaphore := make(chan struct{}, maxRPS)

    // Worker function
    runWorker := func(value string) {
        defer wg.Done()
        defer func() { <-semaphore }()

        atomic.AddInt32(&activeWorkers, 1)
        defer atomic.AddInt32(&activeWorkers, -1)

        // Create the brute-forced request
        bruteforcedReq := createBruteforcedRequest(&baseReq, param, value)

        // Send it
        status, err := sendRequest(bruteforcedReq, ws, db, attackID)
        if err != nil {
            slog.Error("Failed to send bruteforced request", "error", err)
            return
        }

        // Track total requests
        atomic.AddInt32(&totalCount, 1)

        // If success, increment successCount
        if status >= http.StatusOK && status < http.StatusMultipleChoices {
            newCount := atomic.AddInt32(&successCount, 1)
            if int(newCount) >= maxSuccess {
                slog.Info("Max success count reached", "count", maxSuccess, "requests", totalCount)
                once.Do(func() {
                    close(done)
                })
            }
        }
    }

    // Launch workers (one for each dictionary item)
workerLoop:
    for _, val := range dictionary {
        select {
        case <-done:
            break workerLoop
        default:
            wg.Add(1)
            semaphore <- struct{}{}
            go runWorker(val)
        }
    }

    // Wait for all workers to finish
    wg.Wait()

    // Update the final status of the attack
    attack := models.FuzzAttack{}
    db.First(&attack, attackID)
    if successCount > 0 {
        attack.Status = "success"
    } else {
        attack.Status = "completed"
    }
    db.Save(&attack)

    // Notify the UI or WebSocket listeners
    msg := &types.AttackCompleteMessage{
        Type: types.MessageTypeAttackComplete,
    }
    ws.Broadcast <- msg

    slog.Info("Brute force attack completed", "successCount", successCount, "totalCount", totalCount)
    return nil
}

// createBruteforcedRequest modifies the body param with the given 'value'
func createBruteforcedRequest(originalReq *models.Request, key, value string) *http.Request {
    // Parse the original body
    var body map[string]interface{}
    err := json.Unmarshal([]byte(originalReq.Body), &body)
    if err != nil {
        slog.Error("Failed to parse original body", "error", err)
    }

    // If the body was empty or not JSON, initialize
    if body == nil {
        body = make(map[string]interface{})
    }

    // Update the specified key with the brute-forced value
    body[key] = value

    // Convert the updated body back to JSON
    updatedBody, _ := json.Marshal(body)

    // Create a new http.Request with the updated JSON body
    req, _ := http.NewRequest(originalReq.Method, originalReq.URL, strings.NewReader(string(updatedBody)))

    // Re-set the headers from the original request
    headerLines := strings.Split(originalReq.Headers, "\n")
    for _, line := range headerLines {
        line = strings.TrimSpace(line)
        if line != "" {
            parts := strings.SplitN(line, ":", 2)
            if len(parts) == 2 {
                headerKey := strings.TrimSpace(parts[0])
                headerValue := strings.TrimSpace(parts[1])
                req.Header.Add(headerKey, headerValue)
            }
        }
    }

    // Update Content-Length header
    req.Header.Set("Content-Length", fmt.Sprintf("%d", len(updatedBody)))

    return req
}

// sendRequest is identical to the fuzz version, no changes needed
// but I'll copy it here for completeness with zero placeholders

func sendRequest(req *http.Request, ws *websocket.Pool, db *gorm.DB, attackID uint) (int, error) {
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return 0, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    // early return if not 2xx
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

    // Save results in the database
    fuzzResult := &models.FuzzResult{
        FuzzAttackID: attackID,
        Hostname:     req.URL.Hostname(),
        IpAddress:    req.RemoteAddr,
        Port:         req.URL.Port(),
        Scheme:       req.URL.Scheme,
        URL:          req.URL.String(),
        Endpoint:     req.URL.Path,
        Request:      string(requestHeaders) + "\n" + string(requestBody),
        Response:     string(responseBody),
        StatusCode:   resp.StatusCode,
    }
    res := db.Create(fuzzResult)
    if res.Error != nil {
        slog.Error("Failed to save fuzz result", "error", res.Error)
    }

    return resp.StatusCode, nil
}
