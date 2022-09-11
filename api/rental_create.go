package api

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/G1GACHADS/stashable-backend/backend"
	"github.com/G1GACHADS/stashable-backend/core/mime"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slices"
)

const maxImageUploads = 4

var supportedImageTypes = []string{
	"image/jpeg",
	"image/png",
	"image/webp",
}

type RentalCreateParams struct {
	Name         string                     `form:"name"`
	Description  string                     `form:"description"`
	Weight       float64                    `form:"weight"`
	Width        float64                    `form:"width"`
	Height       float64                    `form:"height"`
	Length       float64                    `form:"length"`
	Quantity     int                        `form:"quantity"`
	PaidAnnually bool                       `form:"paid_annually"`
	ShippingType backend.RentalShippingType `form:"shipping_type"`
	CategoryID   int64                      `form:"category_id"`
	RoomID       int64                      `form:"room_id"`
}

func (p RentalCreateParams) Validate() error {
	if err := requiredFields(map[string]any{
		"name":          p.Name,
		"description":   p.Description,
		"weight":        p.Weight,
		"width":         p.Width,
		"height":        p.Height,
		"length":        p.Length,
		"quantity":      p.Quantity,
		"shipping_type": p.ShippingType,
		"category_id":   p.CategoryID,
		"room_id":       p.RoomID,
	}); err != nil {
		return err
	}

	if !slices.Contains([]backend.RentalShippingType{
		backend.RentalSelfServiceShipping,
		backend.RentalDeliveryShipping}, p.ShippingType) {
		return errors.New("invalid rental type (valid => 'self-storage' or 'disposal')")
	}

	return nil
}

func (h *handler) RentalCreate(c *fiber.Ctx) error {
	c.Accepts(fiber.MIMEMultipartForm)

	userID := int64(c.Locals("userID").(float64))
	warehouseID, err := c.ParamsInt("warehouseID")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide a valid warehouse id")
	}

	var params RentalCreateParams
	if err := c.BodyParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := params.Validate(); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if len(form.File["images[]"]) > maxImageUploads {
		return fiber.NewError(fiber.StatusBadRequest,
			fmt.Sprintf("Maximum number of images is %d", maxImageUploads))
	}

	images := make([]backend.RentalCreateMediaInput, len(form.File["images[]"]))

	for idx, image := range form.File["images[]"] {
		file, err := image.Open()
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
		}

		if !mime.Contains(file, supportedImageTypes) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message":         "Invalid image type",
				"file":            image.Filename,
				"supported_types": supportedImageTypes,
			})
		}

		images[idx] = backend.RentalCreateMediaInput{
			File:          file,
			FileHeader:    image,
			FileExtension: filepath.Ext(image.Filename),
		}
	}

	rentalID, err := h.backend.RentalCreate(c.Context(), backend.RentalCreateInput{
		UserID:       userID,
		WarehouseID:  int64(warehouseID),
		CategoryID:   params.CategoryID,
		RoomID:       params.RoomID,
		Images:       images,
		Description:  params.Description,
		PaidAnnually: params.PaidAnnually,
		Name:         params.Name,
		Weight:       params.Weight,
		Width:        params.Width,
		Height:       params.Height,
		Length:       params.Length,
		Quantity:     params.Quantity,
		ShippingType: params.ShippingType,
	})
	if err != nil {
		if errors.Is(err, backend.ErrWarehouseOrCategoryOrRoomDoesNotExists) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		return fiber.NewError(fiber.StatusInternalServerError, "There was a problem on our side")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Successfully booked a rent, please pay the fee to finish the rental process",
		"rental_id":   rentalID,
		"payment_url": fmt.Sprintf("%s/rent/%d", h.appCfg.Address, rentalID),
	})
}
