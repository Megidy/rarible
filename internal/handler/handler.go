package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Megidy/rarible/internal/domain/constants"
	businesserrors "github.com/Megidy/rarible/internal/domain/errors"
	"github.com/Megidy/rarible/internal/domain/model"
	"github.com/Megidy/rarible/internal/handler/dto"
	"github.com/Megidy/rarible/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const (
	idParam = "id"
)

type NFTHandler struct {
	nftService service.NFTService
}

func NewNFTHandler(nftService service.NFTService) *NFTHandler {
	return &NFTHandler{
		nftService: nftService,
	}
}

// GetOwnership godoc
// @Summary Get NFT ownership information
// @Description Retrieves ownership details for a specific NFT by its ID
// @Tags NFT
// @Accept json
// @Produce json
// @Param id path string true "NFT ID" Example(0x06012c8cf97bead5deae237070f9587f8e7a266d:123456)
// @Success 200 {object} dto.GeneralResponse{data=model.OwnershipDTO} "Successfully retrieved ownership data"
// @Failure 400 {object} dto.GeneralResponse "Invalid request parameters"
// @Failure 404 {object} dto.GeneralResponse "NFT ownership not found"
// @Failure 500 {object} dto.GeneralResponse "Internal server error"
// @Router /ownerships/{id} [get]
func (h *NFTHandler) GetOwnership(ctx echo.Context) error {
	id := getFromParam(ctx, idParam)

	ownership, err := h.nftService.GetOwnershipByID(ctx.Request().Context(), id)
	if err != nil {
		msg := "failed to get ownership by id"
		log.Error().Err(err).Msg(msg)

		resp := dto.NewGeneralResponse(nil, constants.StatusFailed, msg, err.Error(), http.StatusInternalServerError)
		switch {
		case errors.Is(err, businesserrors.ErrInvalidRequest):
			resp.Status.StatusCode = http.StatusBadRequest
			return ctx.JSON(http.StatusBadRequest, resp)
		case errors.Is(err, businesserrors.ErrNotFound):
			resp.Status.StatusCode = http.StatusNotFound
			return ctx.JSON(http.StatusNotFound, resp)
		default:
			resp.Status.StatusCode = http.StatusInternalServerError
			return ctx.JSON(http.StatusInternalServerError, resp)
		}
	}

	resp := dto.NewGeneralResponse(ownership, constants.StatusRetrieved, "successfully retrieved data", constants.StrEmpty, http.StatusOK)
	return ctx.JSON(http.StatusOK, resp)
}

// GetTraitRarities godoc
// @Summary Get trait rarities for NFTs
// @Description returns rarity scores for specified traits in an NFT collection
// @Tags NFT
// @Accept json
// @Produce json
// @Param request body model.TraitRarityRequestDTO true "Trait rarity request parameters"
// @Success 200 {object} dto.GeneralResponse{data=model.TraitRarityResponseDTO} "Successfully calculated trait rarities"
// @Failure 400 {object} dto.GeneralResponse "Invalid request body or parameters"
// @Failure 404 {object} dto.GeneralResponse "Collection or traits not found"
// @Failure 500 {object} dto.GeneralResponse "Internal server error"
// @Router /trait-rarities [post]
func (h *NFTHandler) GetTraitRarities(ctx echo.Context) error {
	var req model.TraitRarityRequestDTO

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewGeneralResponse(
			nil,
			constants.StatusFailed,
			"invalid request body", err.Error(), http.StatusBadRequest))
	}

	err = h.validateTraitRarityRequest(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewGeneralResponse(
			nil,
			constants.StatusFailed,
			"invalid request body", err.Error(), http.StatusBadRequest))
	}

	traitRarityResponse, err := h.nftService.GetTraitRarity(ctx.Request().Context(), req)
	if err != nil {
		msg := "failed to get trarity rarity"
		log.Error().Err(err).Msg(msg)

		resp := dto.NewGeneralResponse(nil, constants.StatusFailed, msg, err.Error(), http.StatusInternalServerError)
		switch {
		case errors.Is(err, businesserrors.ErrInvalidRequest):
			resp.Status.StatusCode = http.StatusBadRequest
			return ctx.JSON(http.StatusBadRequest, resp)
		case errors.Is(err, businesserrors.ErrNotFound):
			resp.Status.StatusCode = http.StatusNotFound
			return ctx.JSON(http.StatusNotFound, resp)
		default:
			resp.Status.StatusCode = http.StatusInternalServerError
			return ctx.JSON(http.StatusInternalServerError, resp)
		}
	}

	resp := dto.NewGeneralResponse(traitRarityResponse, constants.StatusRetrieved, "successfully retrieved data", constants.StrEmpty, http.StatusOK)
	return ctx.JSON(http.StatusOK, resp)

}

func (h *NFTHandler) validateTraitRarityRequest(req *model.TraitRarityRequestDTO) error {
	if req.CollectionID == "" {
		return fmt.Errorf("invalid collection id")
	}

	for _, property := range req.Properties {
		if property.Key == "" || property.Value == "" {
			return fmt.Errorf("invalid property param")
		}
	}
	return nil
}
