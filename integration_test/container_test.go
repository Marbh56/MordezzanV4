package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mordezzanV4/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestContainerCRUDIntegration tests the CRUD operations for containers
func TestContainerCRUDIntegration(t *testing.T) {
	log := NewTestLogger(t)
	log.Section("Container CRUD Integration Test")
	log.Step("Setting up test environment")

	testApp, cleanup := setupTestAppWithFullSchema(t)
	defer cleanup()

	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Success("Test server started at %s", server.URL)

	var containerID int64

	t.Run("Create Container", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Creating new container")

		containerData := models.CreateContainerInput{
			Name:         "Backpack",
			MaxWeight:    50,
			AllowedItems: "Books, equipment, small items",
			Cost:         15.0,
		}

		log.Info("Container Name: %s", containerData.Name)
		log.Info("Max Weight: %d lb, Cost: %.1f gp", containerData.MaxWeight, containerData.Cost)

		payload, err := json.Marshal(containerData)
		if !log.CheckNoError(err, "Marshal container data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/containers"
		log.Info("Sending POST request to %s", endpoint)

		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			log.Error("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var createdContainer models.Container
		if err := json.NewDecoder(resp.Body).Decode(&createdContainer); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		containerID = createdContainer.ID
		log.Success("Container created with ID: %d", containerID)

		if createdContainer.Name != containerData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s", containerData.Name, createdContainer.Name)
			t.Errorf("Expected container name %s, got %s", containerData.Name, createdContainer.Name)
		}

		if createdContainer.MaxWeight != containerData.MaxWeight {
			log.Error("Max weight mismatch. Expected: %d, Got: %d", containerData.MaxWeight, createdContainer.MaxWeight)
			t.Errorf("Expected max weight %d, got %d", containerData.MaxWeight, createdContainer.MaxWeight)
		}

		if createdContainer.Cost != containerData.Cost {
			log.Error("Cost mismatch. Expected: %.1f, Got: %.1f", containerData.Cost, createdContainer.Cost)
			t.Errorf("Expected cost %.1f, got %.1f", containerData.Cost, createdContainer.Cost)
		}

		log.Success("Container validation passed")
	})

	log.Separator()

	t.Run("Get Container", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving created container")
		log.Info("Container ID: %d", containerID)

		endpoint := fmt.Sprintf("%s/containers/%d", server.URL, containerID)
		log.Info("Sending GET request to %s", endpoint)

		req, err := http.NewRequest("GET", endpoint, nil)
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var container models.Container
		if err := json.NewDecoder(resp.Body).Decode(&container); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received container: ID=%d, Name=%s", container.ID, container.Name)

		if container.ID != containerID {
			log.Error("Container ID mismatch. Expected: %d, Got: %d", containerID, container.ID)
			t.Errorf("Expected container ID %d, got %d", containerID, container.ID)
		}

		if container.Name != "Backpack" {
			log.Error("Name mismatch. Expected: %s, Got: %s", "Backpack", container.Name)
			t.Errorf("Expected container name 'Backpack', got '%s'", container.Name)
		}

		log.Success("Container data validation passed")
	})

	log.Separator()

	t.Run("List Containers", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving list of containers")

		endpoint := server.URL + "/containers"
		log.Info("Sending GET request to %s", endpoint)

		req, err := http.NewRequest("GET", endpoint, nil)
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var containers []*models.Container
		if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received %d containers in total", len(containers))

		found := false
		for _, c := range containers {
			if c.ID == containerID {
				found = true
				log.Info("Found our test container: ID=%d, Name=%s", c.ID, c.Name)
				break
			}
		}

		if !found {
			log.Error("Container with ID %d not found in containers list", containerID)
			t.Errorf("Container with ID %d not found in containers list", containerID)
		} else {
			log.Success("Container found in the list")
		}
	})

	log.Separator()

	t.Run("Update Container", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Updating container")
		log.Info("Container ID: %d", containerID)

		updateData := models.UpdateContainerInput{
			Name:         "Large Backpack",
			MaxWeight:    75,
			AllowedItems: "Books, equipment, medium items",
			Cost:         25.0,
		}

		log.Info("New container name: %s", updateData.Name)
		log.Info("New Max Weight: %d lb, Cost: %.1f gp", updateData.MaxWeight, updateData.Cost)

		payload, err := json.Marshal(updateData)
		if !log.CheckNoError(err, "Marshal update data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/containers/%d", server.URL, containerID)
		log.Info("Sending PUT request to %s", endpoint)

		req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer(payload))
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var updatedContainer models.Container
		if err := json.NewDecoder(resp.Body).Decode(&updatedContainer); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Updated container data: ID=%d, Name=%s, MaxWeight=%d, Cost=%.1f",
			updatedContainer.ID, updatedContainer.Name, updatedContainer.MaxWeight, updatedContainer.Cost)

		if updatedContainer.Name != updateData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s",
				updateData.Name, updatedContainer.Name)
			t.Errorf("Expected name %s, got %s", updateData.Name, updatedContainer.Name)
		}

		if updatedContainer.MaxWeight != updateData.MaxWeight {
			log.Error("Max weight mismatch. Expected: %d, Got: %d",
				updateData.MaxWeight, updatedContainer.MaxWeight)
			t.Errorf("Expected max weight %d, got %d", updateData.MaxWeight, updatedContainer.MaxWeight)
		}

		if updatedContainer.Cost != updateData.Cost {
			log.Error("Cost mismatch. Expected: %.1f, Got: %.1f",
				updateData.Cost, updatedContainer.Cost)
			t.Errorf("Expected cost %.1f, got %.1f", updateData.Cost, updatedContainer.Cost)
		}

		log.Success("Container update validation passed")
	})

	log.Separator()

	t.Run("Delete Container", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Deleting container")
		log.Info("Container ID: %d", containerID)

		endpoint := fmt.Sprintf("%s/containers/%d", server.URL, containerID)
		log.Info("Sending DELETE request to %s", endpoint)

		req, err := http.NewRequest("DELETE", endpoint, nil)
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			log.Error("Expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)
		log.Info("Verifying container deletion by attempting to retrieve it")
		getReq, err := http.NewRequest("GET", endpoint, nil)
		if !log.CheckNoError(err, "Create verification request") {
			t.Fatal("Test failed")
		}
		getReq.Header.Set("Accept", "application/json")

		getResp, err := http.DefaultClient.Do(getReq)
		if !log.CheckNoError(err, "Send verification request") {
			t.Fatal("Test failed")
		}
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusNotFound {
			log.Error("Expected status %d after deletion, got %d",
				http.StatusNotFound, getResp.StatusCode)
			t.Fatalf("Expected status %d after deletion, got %d",
				http.StatusNotFound, getResp.StatusCode)
		}
		log.Success("Container confirmed deleted (received 404 Not Found)")
	})

	log.Section("CONTAINER CRUD INTEGRATION TEST COMPLETED SUCCESSFULLY")
}
