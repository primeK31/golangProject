package http_delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"golangproject/internal/services/bet"
	"golangproject/internal/services/middleware"
	"golangproject/pkg/reqresp"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sony/gobreaker"
)

func BetHandler(w http.ResponseWriter, r *http.Request, s *bet.SecondServiceClient) {
	var clientData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&clientData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := BetHandlerLogic(r.Context(), clientData, s)
	if err != nil {
		log.Printf("Handler error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set headers
	for key, value := range response.Headers {
		w.Header().Set(key, value)
	}
	w.WriteHeader(response.StatusCode)
	w.Write(response.Body)
}

func BetHandlerLogic(ctx context.Context, clientData map[string]interface{},
	s *bet.SecondServiceClient) (*reqresp.HandlerResponse, error) {

	userIDVal := ctx.Value(middleware.CurrentUserKey)
	if userIDVal == nil {
		return &reqresp.HandlerResponse{
			StatusCode: http.StatusUnauthorized,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       []byte(`{"status":"error","message":"Missing user context"}`),
		}, nil
	}

	// Validate type assertion
	userID := userIDVal
	fmt.Println(userIDVal)
	if userID == uuid.Nil {
		return &reqresp.HandlerResponse{
			StatusCode: http.StatusUnauthorized,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       []byte(`{"status":"error","message":"Invalid user context"}`),
		}, nil
	}
	clientData["userId"] = userID
	// Use userID in your logic
	dataForSecondService := map[string]interface{}{
		"client_data":      clientData,
		"first_service_id": "service-1",
		"timestamp":        time.Now().Unix(),
		"processed_by":     "first-service",
	}
    responseBody, err := s.SendRequest(ctx, dataForSecondService)
    if err != nil {
        log.Printf("Error calling second service: %v", err)

        if err == gobreaker.ErrOpenState {
            response := map[string]string{
                "status":  "error",
                "message": "Second service is currently unavailable, please try again later",
            }
            bodyBytes, marshalErr := json.Marshal(response)
            if marshalErr != nil {
                return nil, fmt.Errorf("failed to marshal JSON response: %w", marshalErr)
            }
            return &reqresp.HandlerResponse{
                StatusCode: http.StatusServiceUnavailable,
                Headers:    map[string]string{"Content-Type": "application/json"},
                Body:       bodyBytes,
            }, nil
        }

        return &reqresp.HandlerResponse{
            StatusCode: http.StatusInternalServerError,
            Headers:    map[string]string{"Content-Type": "text/plain"},
            Body:       []byte("Error processing request"),
        }, nil
    }

    log.Printf("Second service response: %s", responseBody)

    response := map[string]string{
        "status":  "success",
        "message": "Request processed and forwarded to second service",
    }
    bodyBytes, err := json.Marshal(response)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal success response: %w", err)
    }

    return &reqresp.HandlerResponse{
        StatusCode: http.StatusOK,
        Headers:    map[string]string{"Content-Type": "application/json"},
        Body:       bodyBytes,
    }, nil
}
