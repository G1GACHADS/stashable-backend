package backend

import "context"

type RoomGetOutput struct {
	Attributes Room `json:"attributes"`
}

func (b *backend) RoomGet(ctx context.Context, roomID int64) (RoomGetOutput, error) {
	var room RoomGetOutput

	err := b.clients.DB.QueryRow(ctx, "SELECT * FROM rooms WHERE id = $1", roomID).Scan(
		&room.Attributes.ID,
		&room.Attributes.WarehouseID,
		&room.Attributes.ImageURL,
		&room.Attributes.Name,
		&room.Attributes.Width,
		&room.Attributes.Height,
		&room.Attributes.Length,
		&room.Attributes.Price,
	)
	if err != nil {
		return RoomGetOutput{}, err
	}

	return room, nil
}
