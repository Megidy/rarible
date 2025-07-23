package service

import (
	"context"
	"net/http"
	"testing"
	"time"

	client "github.com/Megidy/rarible/internal/client/mock"
	businesserrors "github.com/Megidy/rarible/internal/domain/errors"
	"github.com/Megidy/rarible/internal/domain/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	timeout = time.Second * 10
	id      = "id-123"
)

func TestGetOwnershipByID(t *testing.T) {
	t.Run("ShouldReturnValidValue", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		ownership := &model.OwnershipDTO{
			ID:         id,
			Owner:      "0x123456",
			StatusCode: http.StatusOK,
		}

		client := client.NewMockRaribleClient(ctrl)
		client.EXPECT().GetOwnershipByID(gomock.Any(), id).Return(ownership, nil)

		service := NewNFTService(client)

		resp, err := service.GetOwnershipByID(ctx, id)
		require.NoError(t, err)

		expectedOwnership := *ownership
		actualOwnership := *resp

		require.Equal(t, expectedOwnership, actualOwnership)
	})
	t.Run("ShouldReturnError", func(t *testing.T) {
		t.Run("ShouldReturn400", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			ownership := &model.OwnershipDTO{
				StatusCode: http.StatusBadRequest,
			}
			client := client.NewMockRaribleClient(ctrl)
			client.EXPECT().GetOwnershipByID(gomock.Any(), id).Return(ownership, nil)

			service := NewNFTService(client)

			_, err := service.GetOwnershipByID(ctx, id)

			expectedError := businesserrors.ErrInvalidRequest

			require.Error(t, expectedError, err)
			require.ErrorIs(t, err, expectedError)
		})
		t.Run("ShouldReturn404", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			ownership := &model.OwnershipDTO{
				StatusCode: http.StatusNotFound,
			}
			client := client.NewMockRaribleClient(ctrl)
			client.EXPECT().GetOwnershipByID(gomock.Any(), id).Return(ownership, nil)

			service := NewNFTService(client)

			_, err := service.GetOwnershipByID(ctx, id)

			expectedError := businesserrors.ErrNotFound

			require.Error(t, expectedError, err)
			require.ErrorIs(t, err, expectedError)
		})
		t.Run("ShouldReturn500", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			ownership := &model.OwnershipDTO{
				StatusCode: http.StatusInternalServerError,
			}
			client := client.NewMockRaribleClient(ctrl)
			client.EXPECT().GetOwnershipByID(gomock.Any(), id).Return(ownership, nil)

			service := NewNFTService(client)

			_, err := service.GetOwnershipByID(ctx, id)

			expectedError := businesserrors.ErrSomethingWentWrong

			require.Error(t, expectedError, err)
			require.ErrorIs(t, err, expectedError)
		})
	})
}

func TestGetTraitRarity(t *testing.T) {
	t.Run("ShouldReturnValidValue", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		mockRequest := model.TraitRarityRequestDTO{
			CollectionID: "ETHEREUM:0x123",
			Properties: []model.TraitPropertyInput{
				{Key: "Hat", Value: "Halo"},
			},
		}

		mockResponse := &model.TraitRarityResponseDTO{
			Continuation: "token123",
			Traits: []model.ExtendedTraitProperty{
				{Key: "Hat", Value: "Halo", Rarity: "1.2"},
			},
			StatusCode: http.StatusOK,
		}

		client := client.NewMockRaribleClient(ctrl)
		client.EXPECT().GetTraitRarity(gomock.Any(), &mockRequest).Return(mockResponse, nil)

		service := NewNFTService(client)

		resp, err := service.GetTraitRarity(ctx, mockRequest)
		require.NoError(t, err)

		require.Equal(t, mockResponse.Continuation, resp.Continuation)
		require.Equal(t, len(mockResponse.Traits), len(resp.Traits))
		require.Equal(t, mockResponse.Traits[0].Key, resp.Traits[0].Key)
		require.Equal(t, mockResponse.Traits[0].Value, resp.Traits[0].Value)
		require.Equal(t, mockResponse.Traits[0].Rarity, resp.Traits[0].Rarity)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("ShouldReturnError", func(t *testing.T) {
		t.Run("ShouldReturn400", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			mockRequest := model.TraitRarityRequestDTO{
				CollectionID: "ETHEREUM:0x123",
				Properties: []model.TraitPropertyInput{
					{Key: "Hat", Value: "Halo"},
				},
			}

			mockResponse := &model.TraitRarityResponseDTO{
				StatusCode: http.StatusBadRequest,
				Message:    "bad request",
			}

			client := client.NewMockRaribleClient(ctrl)
			client.EXPECT().GetTraitRarity(gomock.Any(), &mockRequest).Return(mockResponse, nil)

			service := NewNFTService(client)

			_, err := service.GetTraitRarity(ctx, mockRequest)
			expectedError := businesserrors.ErrInvalidRequest

			require.Error(t, err)
			require.ErrorIs(t, err, expectedError)
		})

		t.Run("ShouldReturn404", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			mockRequest := model.TraitRarityRequestDTO{
				CollectionID: "ETHEREUM:0x123",
				Properties: []model.TraitPropertyInput{
					{Key: "Hat", Value: "Halo"},
				},
			}

			mockResponse := &model.TraitRarityResponseDTO{
				StatusCode: http.StatusNotFound,
				Message:    "not found",
			}

			client := client.NewMockRaribleClient(ctrl)
			client.EXPECT().GetTraitRarity(gomock.Any(), &mockRequest).Return(mockResponse, nil)

			service := NewNFTService(client)

			_, err := service.GetTraitRarity(ctx, mockRequest)
			expectedError := businesserrors.ErrNotFound

			require.Error(t, err)
			require.ErrorIs(t, err, expectedError)
		})

		t.Run("ShouldReturn500", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			mockRequest := model.TraitRarityRequestDTO{
				CollectionID: "ETHEREUM:0x123",
				Properties: []model.TraitPropertyInput{
					{Key: "Hat", Value: "Halo"},
				},
			}

			mockResponse := &model.TraitRarityResponseDTO{
				StatusCode: http.StatusInternalServerError,
				Message:    "internal error",
			}

			client := client.NewMockRaribleClient(ctrl)
			client.EXPECT().GetTraitRarity(gomock.Any(), &mockRequest).Return(mockResponse, nil)

			service := NewNFTService(client)

			_, err := service.GetTraitRarity(ctx, mockRequest)
			expectedError := businesserrors.ErrSomethingWentWrong

			require.Error(t, err)
			require.ErrorIs(t, err, expectedError)
		})
	})
}
