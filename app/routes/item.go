package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/TrevorEdris/api-template/app/controller"
	"github.com/TrevorEdris/api-template/app/model/item"
)

type (
	// Item is the controller for the Item-based routes.
	Item struct {
		controller.Controller
	}

	// ItemPostRequest defines the HTTP request body for the POST handler.
	itemPostRequest struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	// ItemPutRequest defines the HTTP request body for the PUT handler.
	itemPutRequest struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	// ItemJSONResponse defines the HTTP response body of a request which returns an item.Model.
	itemJSONResponse struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}
)

func (c Item) Get(ctx echo.Context) error {
	ctx.Logger().Debug("Processing GET request")
	resp := controller.NewJSONResponse(ctx)
	resp.StatusCode = http.StatusOK

	id := ctx.Param("id")
	it, err := c.Container.ItemRepo.Get(ctx.Request().Context(), id)
	if err != nil {
		switch err {
		case item.ErrItemNotFound:
			return c.RenderErrorResponse(ctx, http.StatusNotFound, err, "item (id="+id+") not found")
		default:
			return c.RenderErrorResponse(ctx, http.StatusInternalServerError, err, "failed to get item")
		}
	}

	resp.Body = modelToItemJSONResponse(it)

	return c.RenderJSONResponse(ctx, resp)
}

func (c Item) Post(ctx echo.Context) error {
	ctx.Logger().Debug("Processing POST request")
	resp := controller.NewJSONResponse(ctx)
	resp.StatusCode = http.StatusCreated

	ipr := new(itemPostRequest)
	err := ctx.Bind(ipr)
	if err != nil {
		return c.RenderErrorResponse(ctx, http.StatusBadRequest, err, "invalid request body")
	}

	err = ctx.Validate(ipr)
	if err != nil {
		return c.RenderErrorResponse(ctx, http.StatusBadRequest, err, "request body missing required fields")
	}

	it, err := c.Container.ItemRepo.Create(ctx.Request().Context(), itemPostRequestToModel(ipr))
	if err != nil {
		return c.RenderErrorResponse(ctx, http.StatusInternalServerError, err, "failed to create item")
	}

	resp.Body = modelToItemJSONResponse(it)

	return c.RenderJSONResponse(ctx, resp)
}

func (c Item) Put(ctx echo.Context) error {
	ctx.Logger().Debug("Processing PUT request")
	resp := controller.NewJSONResponse(ctx)
	resp.StatusCode = http.StatusCreated

	id := ctx.Param("id")
	ipr := new(itemPutRequest)
	err := ctx.Bind(ipr)
	if err != nil {
		return c.RenderErrorResponse(ctx, http.StatusBadRequest, err, "invalid request body")
	}

	err = ctx.Validate(ipr)
	if err != nil {
		return c.RenderErrorResponse(ctx, http.StatusBadRequest, err, "request body missing required fields")
	}

	it, err := c.Container.ItemRepo.Update(ctx.Request().Context(), id, itemPutRequestToModel(ipr))
	if err != nil {
		return c.RenderErrorResponse(ctx, http.StatusInternalServerError, err, "failed to update item")
	}

	resp.Body = modelToItemJSONResponse(it)

	return c.RenderJSONResponse(ctx, resp)
}

func (c Item) Delete(ctx echo.Context) error {
	ctx.Logger().Debug("Processing DELETE request")
	resp := controller.NewJSONResponse(ctx)
	resp.StatusCode = http.StatusNoContent

	id := ctx.Param("id")
	err := c.Container.ItemRepo.Delete(ctx.Request().Context(), id)
	if err != nil {
		switch err {
		case item.ErrItemNotFound:
			return c.RenderErrorResponse(ctx, http.StatusNotFound, err, "item (id="+id+") not found")
		default:
			return c.RenderErrorResponse(ctx, http.StatusInternalServerError, err, "failed to delete item")
		}
	}

	return c.RenderJSONResponse(ctx, resp)
}

func itemPostRequestToModel(ipr *itemPostRequest) item.Model {
	return item.New("", ipr.Name, ipr.Description, ipr.Price)
}

func itemPutRequestToModel(ipr *itemPutRequest) item.Model {
	return item.New("", ipr.Name, ipr.Description, ipr.Price)
}

func modelToItemJSONResponse(m item.Model) itemJSONResponse {
	return itemJSONResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Price:       m.Price,
	}
}
