package client

import (
	"context"

	"github.com/Megidy/rarible/internal/domain/model"
)

type RaribleClient interface {
	// GetOwnershipByID fetches ownership data by ID
	GetOwnershipByID(ctx context.Context, id string) (*model.OwnershipDTO, error)
	// GetTraitRarity returns rarity of a given trait
	GetTraitRarity(ctx context.Context, req *model.TraitRarityRequestDTO) (*model.TraitRarityResponseDTO, error)
}
