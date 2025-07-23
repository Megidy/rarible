package service

import (
	"context"

	"github.com/Megidy/rarible/internal/domain/model"
)

type NFTService interface {
	GetOwnershipByID(ctx context.Context, id string) (*model.OwnershipDTO, error)
	GetTraitRarity(ctx context.Context, req model.TraitRarityRequestDTO) (*model.TraitRarityResponseDTO, error)
}
