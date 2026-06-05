package test

import (
	"log"
	"net/http"

	"github.com/mohit838/olario-platform-backend/cmd/internal/dto"
	"github.com/mohit838/olario-platform-backend/cmd/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetTests(w http.ResponseWriter, r *http.Request) {
	// TEST from service layer
	tests, err := h.service.ListTests(r.Context())
	if err != nil {
		log.Println("failed to list tests:", err)
		utils.ApiResponse(w, http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to fetch tests",
			Error:   "internal server error",
		})
		return
	}

	utils.ApiResponse(w, http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "tests fetched successfully",
		Data:    tests,
	})
}

func (h *Handler) CreateTest(w http.ResponseWriter, r *http.Request) {
	var req dto.TestRequest
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		log.Println("failed to decode request body:", err)
		utils.ApiResponse(w, http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "invalid request body",
			Error:   "bad request",
		})
		return
	}

	test, err := h.service.CreateTest(r.Context(), req)
	if err != nil {
		log.Println("failed to create test:", err)
		utils.ApiResponse(w, http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "failed to create test",
			Error:   "internal server error",
		})
		return
	}

	utils.ApiResponse(w, http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "test created successfully",
		Data:    test,
	})
}
