package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Megidy/rarible/internal/domain/constants"
	businesserrors "github.com/Megidy/rarible/internal/domain/errors"
	"github.com/Megidy/rarible/internal/domain/model"
	"github.com/Megidy/rarible/internal/handler/dto"
	service "github.com/Megidy/rarible/internal/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestNFTHandler_GetOwnership(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockNFTService(ctrl)
	h := NewNFTHandler(mockService)

	t.Run("Success", func(t *testing.T) {
		mockOwnership := &model.OwnershipDTO{ID: "id-123", Owner: "0xabc", StatusCode: http.StatusOK}
		mockService.EXPECT().GetOwnershipByID(gomock.Any(), "id-123").Return(mockOwnership, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/ownerships/id-123", http.NoBody)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("id-123")

		err := h.GetOwnership(c)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)

		var resp dto.GeneralResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Equal(t, constants.StatusRetrieved, resp.Status.Status)
	})

	t.Run("NotFoundError", func(t *testing.T) {
		mockService.EXPECT().GetOwnershipByID(gomock.Any(), "missing-id").Return(nil, businesserrors.ErrNotFound)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/ownerships/missing-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("missing-id")

		err := h.GetOwnership(c)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("InvalidRequestError", func(t *testing.T) {
		mockService.EXPECT().GetOwnershipByID(gomock.Any(), "bad-id").Return(nil, businesserrors.ErrInvalidRequest)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/ownerships/bad-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("bad-id")

		err := h.GetOwnership(c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockService.EXPECT().GetOwnershipByID(gomock.Any(), "error-id").Return(nil, errors.New("some error"))

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/ownerships/error-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("error-id")

		err := h.GetOwnership(c)
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestNFTHandler_GetTraitRarities(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockNFTService(ctrl)
	h := NewNFTHandler(mockService)

	t.Run("Success", func(t *testing.T) {
		reqBody := model.TraitRarityRequestDTO{
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

		mockService.EXPECT().GetTraitRarity(gomock.Any(), reqBody).Return(mockResponse, nil)

		e := echo.New()
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/traits/rarity", strings.NewReader(string(bodyBytes)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.GetTraitRarities(c)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)

		var resp dto.GeneralResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Equal(t, constants.StatusRetrieved, resp.Status.Status)
	})

	t.Run("BadRequestInvalidBody", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/traits/rarity", bytes.NewReader([]byte("invalid json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.GetTraitRarities(c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("BadRequestInvalidFields", func(t *testing.T) {
		e := echo.New()
		invalidReq := model.TraitRarityRequestDTO{
			CollectionID: "",
			Properties:   []model.TraitPropertyInput{{Key: "", Value: ""}},
		}
		bodyBytes, _ := json.Marshal(invalidReq)
		req := httptest.NewRequest(http.MethodPost, "/traits/rarity", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.GetTraitRarities(c)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		reqBody := model.TraitRarityRequestDTO{
			CollectionID: "ETHEREUM:0x123",
			Properties: []model.TraitPropertyInput{
				{Key: "Hat", Value: "Halo"},
			},
		}

		mockService.EXPECT().GetTraitRarity(gomock.Any(), reqBody).Return(nil, errors.New("some error"))

		e := echo.New()
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/traits/rarity", bytes.NewReader(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.GetTraitRarities(c)
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
