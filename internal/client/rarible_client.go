package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Megidy/rarible/internal/domain/model"
)

const (
	httpTimeout = 10 * time.Second
)

type raribleClient struct {
	baseRaribleUrl string
	apiKey         string
	client         *http.Client
}

func NewRaribleClient(apiKey string, baseRaribleUrl string) RaribleClient {
	return &raribleClient{
		baseRaribleUrl: baseRaribleUrl,
		apiKey:         apiKey,
		client: &http.Client{
			Timeout: httpTimeout,
		},
	}
}

// GetOwnershipByID fetches ownership data by ID
func (c *raribleClient) GetOwnershipByID(ctx context.Context, id string) (*model.OwnershipDTO, error) {
	url := fmt.Sprintf("%s/ownerships/%s", c.baseRaribleUrl, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setRequiredHeaders(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var ownership model.OwnershipDTO
	if err := json.NewDecoder(resp.Body).Decode(&ownership); err != nil {
		return nil, fmt.Errorf("failed to decode ownership response: %w", err)
	}

	ownership.StatusCode = resp.StatusCode

	return &ownership, nil
}

// GetTraitRarity returns rarity of a given trait
func (c *raribleClient) GetTraitRarity(ctx context.Context, dto *model.TraitRarityRequestDTO) (*model.TraitRarityResponseDTO, error) {
	url := fmt.Sprintf("%s/items/traits/rarity", c.baseRaribleUrl)

	bodyBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setRequiredHeaders(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var traitRarity model.TraitRarityResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&traitRarity); err != nil {
		return nil, fmt.Errorf("failed to decode trait rarity response: %w", err)
	}

	traitRarity.StatusCode = resp.StatusCode
	return &traitRarity, nil
}

func (c *raribleClient) setRequiredHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Set("X-API-KEY", c.apiKey)
}
