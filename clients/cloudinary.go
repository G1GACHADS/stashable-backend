package clients

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Cloudinary struct {
	cld    *cloudinary.Cloudinary
	folder string
}

func NewCloudinary(cloud, key, secret, folder string) (*Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(cloud, key, secret)
	if err != nil {
		return nil, err
	}

	return &Cloudinary{cld: cld, folder: folder}, nil
}

func (c *Cloudinary) Upload(ctx context.Context, file any, key string) (string, error) {
	result, err := c.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:   c.folder,
		PublicID: key,
	})
	if err != nil {
		return "", err
	}

	return result.SecureURL, nil
}

func (c *Cloudinary) Delete(ctx context.Context, key string) error {
	_, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: key,
	})
	if err != nil {
		return err
	}

	return nil
}
