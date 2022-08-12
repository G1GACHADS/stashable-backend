package backend

import "context"

func (b backend) CreateCategory(ctx context.Context, name string) (Category, error) {
	category := Category{Name: name}

	err := b.clients.DB.QueryRow(ctx, "INSERT INTO categories (name) VALUES ($1) RETURNING id", name).Scan(&category.ID)
	if err != nil {
		return Category{}, err
	}

	return category, nil
}

func (b backend) DeleteCategory(ctx context.Context, categoryID int64) error {
	var exists bool
	err := b.clients.DB.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM categories WHERE id = $1)", categoryID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return ErrCategoryDoesNotExists
	}

	_, err = b.clients.DB.Exec(ctx, "DELETE FROM categories WHERE id = $1", categoryID)
	return err
}
