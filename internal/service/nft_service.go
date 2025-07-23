package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Megidy/rarible/internal/client"
	businesserrors "github.com/Megidy/rarible/internal/domain/errors"
	"github.com/Megidy/rarible/internal/domain/model"
)

type nftService struct {
	raribleClient client.RaribleClient
}

func NewNFTService(raribleClient client.RaribleClient) NFTService {
	return &nftService{
		raribleClient: raribleClient,
	}
}

func (s *nftService) GetOwnershipByID(ctx context.Context, id string) (*model.OwnershipDTO, error) {
	ownership, err := s.raribleClient.GetOwnershipByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from api: %w", err)
	}
	if ownership.StatusCode != http.StatusOK {
		return nil, s.handleErrors(ownership.StatusCode, ownership.Message)
	}

	// if status code is valid, returning value
	return ownership, nil
}

func (s *nftService) GetTraitRarity(ctx context.Context, req model.TraitRarityRequestDTO) (*model.TraitRarityResponseDTO, error) {
	resp, err := s.raribleClient.GetTraitRarity(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from api: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, s.handleErrors(resp.StatusCode, resp.Message)
	}
	return resp, nil
}

// handleErrors function that maps API statuses to business errors for consistent error handling
func (s *nftService) handleErrors(statusCode int, message string) error {
	switch statusCode {
	case http.StatusBadRequest:
		return fmt.Errorf("%w: %s", businesserrors.ErrInvalidRequest, message)
	case http.StatusNotFound:
		return fmt.Errorf("%w: %s", businesserrors.ErrNotFound, message)
	default:
		return fmt.Errorf("%w: %s", businesserrors.ErrSomethingWentWrong, message)
	}
}
