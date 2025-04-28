package bet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)


type SecondServiceClient struct {
	apiGatewayURL string
	servicePath   string
	circuitBreaker *gobreaker.CircuitBreaker
	httpClient    *http.Client
}

func NewSecondServiceClient(apiGatewayURL, servicePath string) *SecondServiceClient {
	// configure circuit breaker
	cbSettings := gobreaker.Settings{
		Name:        "second-service-cb",
		MaxRequests: 1,
		Interval:    time.Minute,
		Timeout:     time.Minute * 2,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio >= 0.5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit breaker '%s' state changed from %s to %s", name, from, to)
		},
	}

	return &SecondServiceClient{
		apiGatewayURL:  apiGatewayURL,
		servicePath:    servicePath,
		circuitBreaker: gobreaker.NewCircuitBreaker(cbSettings),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *SecondServiceClient) SendRequest(ctx context.Context, data interface{}) (string, error) {
	result, err := s.circuitBreaker.Execute(func() (interface{}, error) {
		return s.doSendRequest(ctx, data)
	})

	if err != nil {
		return "", err
	}

	return result.(string), nil
}


func (s *SecondServiceClient) doSendRequest(ctx context.Context, data interface{}) (interface{}, error) {
	payload, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data is not a valid map")
	}

	clientDataRaw, ok := payload["client_data"]
	if !ok {
		return nil, fmt.Errorf("client_data field missing")
	}

	clientData, ok := clientDataRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("client_data is not a valid map")
	}

	fmt.Println("clientData:", clientData)

	jsonData, err := json.Marshal(clientData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling clientData: %w", err)
	}

	url := s.apiGatewayURL + s.servicePath
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer your-auth-token")

	if correlationID, ok := ctx.Value("correlation_id").(string); ok && correlationID != "" {
		req.Header.Set("X-Correlation-ID", correlationID)
	} else {
		req.Header.Set("X-Correlation-ID", fmt.Sprintf("req-%d", time.Now().UnixNano()))
	}

	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("service returned error status %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}
