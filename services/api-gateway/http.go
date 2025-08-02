package main

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validation
	if reqBody.UserID == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	// jsonBody, _ := json.Marshal(reqBody)
	// reader := bytes.NewReader(jsonBody)

	// Why we need to create a new client for each connection:
	// because if a service is down, we don't want to block the whole application
	// so we created a new client for each connection
	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	// Don't forget to close the client to avoid resource leaks!
	defer tripService.Close()

	tripPreview, err := tripService.Client.PreviewTrip(r.Context(), reqBody.ToProto())
	if err != nil {
		log.Printf("Failed to preview a trip: %v", err)
		http.Error(w, "Failed to preview trip", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: tripPreview}
	writeJSON(w, http.StatusCreated, response)
}
