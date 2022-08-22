package api

import (
	"errors"
	"fmt"
	"os"

	"github.com/G1GACHADS/stashable-backend/backend"
	"github.com/G1GACHADS/stashable-backend/core/logger"
	"github.com/G1GACHADS/stashable-backend/core/mime"
	"github.com/G1GACHADS/stashable-backend/core/nanoid"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slices"
)

const maxImageUploads = 4

var supportedImageTypes = []string{
	"image/avif",
	"image/bmp",
	"image/jpeg",
	"image/png",
	"image/webp",
}

type CreateRentalParams struct {
	Name         string             `form:"name"`
	Description  string             `form:"description"`
	Weight       float64            `form:"weight"`
	Width        float64            `form:"width"`
	Height       float64            `form:"height"`
	Length       float64            `form:"length"`
	Quantity     int                `form:"quantity"`
	PaidAnnually bool               `form:"paid_annually"`
	Type         backend.RentalType `form:"type"`
	CategoryID   int64              `form:"category_id"`
	RoomID       int64              `form:"room_id"`
}

func (p CreateRentalParams) Validate() error {
	if err := requiredFields(map[string]any{
		"name":        p.Name,
		"description": p.Description,
		"weight":      p.Weight,
		"width":       p.Width,
		"height":      p.Height,
		"length":      p.Length,
		"quantity":    p.Quantity,
		"type":        p.Type,
		"category_id": p.CategoryID,
		"room_id":     p.RoomID,
	}); err != nil {
		return err
	}

	if !slices.Contains([]backend.RentalType{
		backend.RentalSelfStorage,
		backend.RentalDisposal}, p.Type) {
		return errors.New("invalid rental type (valid => 'self-storage' or 'disposal')")
	}

	return nil
}

func (h *handler) CreateRental(c *fiber.Ctx) error {
	c.Accepts(fiber.MIMEMultipartForm)

	userID := int64(c.Locals("userID").(float64))
	warehouseID, err := c.ParamsInt("warehouseID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Please provide a valid warehouse id",
		})
	}

	var params CreateRentalParams
	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := params.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if len(form.File["images"]) > maxImageUploads {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": fmt.Sprintf("Maximum number of images is %d", maxImageUploads),
		})
	}

	imageURLs := make([]string, len(form.File["images"]))

	for idx, image := range form.File["images"] {
		file, err := image.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request body",
				"err":     err.Error(),
			})
		}

		if !mime.Contains(file, supportedImageTypes) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message":         "Invalid image type",
				"file":            image.Filename,
				"supported_types": supportedImageTypes,
			})
		}

		fileName := nanoid.Next()
		filePath := h.appCfg.UploadsPath + "/" + fileName + "." + image.Filename
		if err := c.SaveFile(image, filePath); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		imageURLs[idx] = h.appCfg.Address + "/" + filePath
	}

	rentalID, err := h.backend.CreateRental(c.Context(), backend.CreateRentalInput{
		UserID:       userID,
		WarehouseID:  int64(warehouseID),
		CategoryID:   params.CategoryID,
		RoomID:       params.RoomID,
		ImageURLs:    imageURLs,
		Description:  params.Description,
		PaidAnnually: params.PaidAnnually,
		Name:         params.Name,
		Weight:       params.Weight,
		Width:        params.Width,
		Height:       params.Height,
		Length:       params.Length,
		Quantity:     params.Quantity,
		Type:         params.Type,
	})
	if err != nil {
		go func(imageURLs []string) {
			for _, imageURL := range imageURLs {
				path := imageURL[len(h.appCfg.Address):]
				wd, err := os.Getwd()
				if err != nil {
					logger.M.Warnf("failed removing %s: %v", imageURL, err)
					continue
				}

				if err := os.Remove(wd + path); err != nil {
					logger.M.Warnf("failed removing %s: %v", imageURL, err)
				}
			}
		}(imageURLs)

		if errors.Is(err, backend.ErrWarehouseOrCategoryOrRoomDoesNotExists) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "There was a problem on our side",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Successfully booked a rent, please pay the fee to finish the rental process",
		"rental_id":   rentalID,
		"payment_url": fmt.Sprintf("%s/rent/%d", h.appCfg.Address, rentalID),
	})
}
