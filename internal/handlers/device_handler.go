package handlers

import (
	"github.com/gin-gonic/gin"
	"go-backend/internal/dto"
	"go-backend/internal/models"
	"go-backend/internal/services"
	apperror "go-backend/pkg/error"
	"net/http"
	"strconv"
)

type DeviceHandler struct{ svc *services.DeviceService }

func NewDeviceHandler(s *services.DeviceService) *DeviceHandler { return &DeviceHandler{svc: s} }

func (h *DeviceHandler) Create(c *gin.Context) {
	var req dto.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperror.JSONError(c, http.StatusBadRequest, "validation_error", "invalid request payload", err.Error())
		return
	}
	d := models.Device{Name: req.Name, Brand: req.Brand, State: models.State(req.State), CreatedAt: models.NowFormattedTime()}
	id, err := h.svc.Create(c, &d)
	if err != nil {
		httpError(c, err)
		return
	}
	d.ID = id
	c.JSON(http.StatusCreated, dto.FromModel(&d))
}

func (h *DeviceHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	d, err := h.svc.Get(c, id)
	if err != nil {
		httpError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.FromModel(d))
}

func (h *DeviceHandler) List(c *gin.Context) {
	brand := c.Query("brand")
	state := c.Query("state")
	list, err := h.svc.List(c, brand, state)
	if err != nil {
		httpError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.FromModels(list))
}

func (h *DeviceHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req dto.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperror.JSONError(c, http.StatusBadRequest, "validation_error", "invalid request payload", err.Error())
		return
	}
	var created models.FormattedTime
	if req.CreatedAt != nil {
		created = models.NewFormattedTime(*req.CreatedAt)
	} else {
		ex, err := h.svc.Get(c, id)
		if err != nil {
			httpError(c, err)
			return
		}
		created = ex.CreatedAt
	}
	d := models.Device{Name: req.Name, Brand: req.Brand, State: models.State(req.State), CreatedAt: created}
	if err := h.svc.Update(c, id, &d); err != nil {
		httpError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *DeviceHandler) Patch(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req dto.PatchDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperror.JSONError(c, http.StatusBadRequest, "validation_error", "invalid request payload", err.Error())
		return
	}
	m := map[string]any{}
	if req.Name != nil {
		m["name"] = *req.Name
	}
	if req.Brand != nil {
		m["brand"] = *req.Brand
	}
	if req.State != nil {
		m["state"] = *req.State
	}
	if err := h.svc.Patch(c, id, m); err != nil {
		httpError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *DeviceHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c, id); err != nil {
		httpError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func parseID(c *gin.Context) (int64, bool) {
	sid := c.Param("id")
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return id, true
}

func httpError(c *gin.Context, err error) {
	switch err {
	case models.ErrCannotDeleteInUse:
		apperror.JSONError(c, http.StatusConflict, "in_use_delete_blocked", err.Error(), nil)
	case models.ErrCannotUpdateCreated:
		apperror.JSONError(c, http.StatusUnprocessableEntity, "cannot_update_created_at", err.Error(), nil)
	case models.ErrCannotUpdateFields:
		apperror.JSONError(c, http.StatusUnprocessableEntity, "cannot_update_name_brand_in_use", err.Error(), nil)
	case models.ErrInvalidState:
		apperror.JSONError(c, http.StatusUnprocessableEntity, "invalid_state", err.Error(), nil)
	default:
		apperror.JSONError(c, http.StatusInternalServerError, "internal_error", err.Error(), nil)
	}
}
