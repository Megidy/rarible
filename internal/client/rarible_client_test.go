package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Megidy/rarible/internal/domain/model"
	"github.com/stretchr/testify/require"
)

func TestGetOwnershipByID(t *testing.T) {
	t.Run("ShouldPass_DefaultCase(mocked_200_response_from_server)", func(t *testing.T) {
		mockOwnership := model.OwnershipDTO{
			ID:    "test-id",
			Owner: "owner-address",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/ownerships/test-id", r.URL.Path)
			require.Equal(t, "application/json", r.Header.Get("Accept"))
			require.Equal(t, "application/json", r.Header.Get("content-type"))
			require.Equal(t, "test-api-key", r.Header.Get("X-API-KEY"))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockOwnership)
		}))
		defer server.Close()

		client := NewRaribleClient("test-api-key", server.URL)

		ownership, err := client.GetOwnershipByID(context.Background(), "test-id")
		require.NoError(t, err)

		expectedId := mockOwnership.ID
		actualId := ownership.ID

		expectedStatusCode := http.StatusOK
		actualStatusCode := ownership.StatusCode

		expectedOwner := mockOwnership.Owner
		actualOwner := ownership.Owner

		require.Equal(t, expectedId, actualId)
		require.Equal(t, expectedOwner, actualOwner)
		require.Equal(t, expectedStatusCode, actualStatusCode)
	})
	t.Run("ShouldPass_WhenNoOwnershipFound(mocked_404_response_from_server)", func(t *testing.T) {
		mockOwnership := model.OwnershipDTO{
			ID:    "test-id",
			Owner: "owner-address",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/ownerships/test-id" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(mockOwnership)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(mockOwnership)
			}
		}))

		defer server.Close()

		client := NewRaribleClient("test-api-key", server.URL)

		ownership, err := client.GetOwnershipByID(context.Background(), "invalid-id")
		require.NoError(t, err)

		expectedStatusCode := http.StatusNotFound
		actualStatusCode := ownership.StatusCode

		require.Equal(t, expectedStatusCode, actualStatusCode)
	})
	t.Run("ShouldPass_WhenBadRequest(mocked_400_response_from_server)", func(t *testing.T) {
		mockOwnership := model.OwnershipDTO{
			ID:    "test-id",
			Owner: "owner-address",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/ownerships/test-id" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(mockOwnership)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(mockOwnership)
			}
		}))

		defer server.Close()

		client := NewRaribleClient("test-api-key", server.URL)

		ownership, err := client.GetOwnershipByID(context.Background(), "invalid-id")
		require.NoError(t, err)

		expectedStatusCode := http.StatusBadRequest
		actualStatusCode := ownership.StatusCode

		require.Equal(t, expectedStatusCode, actualStatusCode)
	})

}

func TestGetTraitRarity(t *testing.T) {
	t.Run("ShouldPass_DefaultCase(mocked_200_response_from_server)", func(t *testing.T) {
		mockRequest := model.TraitRarityRequestDTO{
			CollectionID: "ETHEREUM:0x60e4d786628fea6478f785a6d7e704777c86a7c6",
			Properties: []model.TraitPropertyInput{
				{Key: "Hat", Value: "Halo"},
			},
		}

		mockResponse := model.TraitRarityResponseDTO{
			Continuation: "token123",
			Traits: []model.ExtendedTraitProperty{
				{
					Key:    "Hat",
					Value:  "Halo",
					Rarity: "1.2",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/items/traits/rarity", r.URL.Path)

			require.Equal(t, "application/json", r.Header.Get("Accept"))
			require.Equal(t, "application/json", r.Header.Get("content-type"))
			require.Equal(t, "test-api-key", r.Header.Get("X-API-KEY"))

			var reqBody model.TraitRarityRequestDTO
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			require.NoError(t, err)
			require.Equal(t, mockRequest.CollectionID, reqBody.CollectionID)
			require.Len(t, reqBody.Properties, 1)
			require.Equal(t, mockRequest.Properties[0], reqBody.Properties[0])

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockResponse)
		}))
		defer server.Close()

		client := NewRaribleClient("test-api-key", server.URL)

		resp, err := client.GetTraitRarity(context.Background(), &mockRequest)
		require.NoError(t, err)
		require.NotNil(t, resp)

		require.Equal(t, mockResponse.Continuation, resp.Continuation)
		require.Equal(t, len(mockResponse.Traits), len(resp.Traits))
		require.Equal(t, mockResponse.Traits[0].Key, resp.Traits[0].Key)
		require.Equal(t, mockResponse.Traits[0].Value, resp.Traits[0].Value)
		require.Equal(t, mockResponse.Traits[0].Rarity, resp.Traits[0].Rarity)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("ShouldPass_WhenBadRequest(mocked_400_response_from_server)", func(t *testing.T) {
		mockRequest := model.TraitRarityRequestDTO{
			CollectionID: "ETHEREUM:0x60e4d786628fea6478f785a6d7e704777c86a7c6",
			Properties: []model.TraitPropertyInput{
				{Key: "Hat", Value: "Halo"},
			},
		}

		mockResponse := model.TraitRarityResponseDTO{
			Continuation: "token123",
			Traits: []model.ExtendedTraitProperty{
				{
					Key:    "Hat",
					Value:  "Halo",
					Rarity: "1.2",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(mockResponse)
		}))
		defer server.Close()

		client := NewRaribleClient("test-api-key", server.URL)

		resp, err := client.GetTraitRarity(context.Background(), &mockRequest)
		require.NoError(t, err)

		expectedStatusCode := http.StatusBadRequest
		actualStatusCode := resp.StatusCode

		require.Equal(t, expectedStatusCode, actualStatusCode)
	})
}
