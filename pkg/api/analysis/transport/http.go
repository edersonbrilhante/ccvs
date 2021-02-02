package transport

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/edersonbrilhante/vilicus"
	"github.com/edersonbrilhante/vilicus/pkg/api/analysis"
)

// Analysis create request
type createReq struct {
	Image string `json:"image" validate:"required"`
}

// NewHTTP creates new analysis http service
func NewHTTP(svc analysis.Service, r *echo.Group) {
	h := HTTP{svc}
	ur := r.Group("/analysis")

	ur.POST("", h.create)
	ur.GET("/:id", h.view)
}

// HTTP represents analysis http service
type HTTP struct {
	svc analysis.Service
}

func (h HTTP) create(c echo.Context) error {
	req := new(createReq)

	if err := c.Bind(req); err != nil {
		return err
	}

	al, err := h.svc.Create(c, vilicus.Analysis{
		Image: req.Image,
	})

	alu := al
	go h.svc.Update(c, alu)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, al)
}

func (h HTTP) view(c echo.Context) error {
	id := c.Param("id")

	result, err := h.svc.View(c, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
