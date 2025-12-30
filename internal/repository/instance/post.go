package instance

import "errors"

func (d *Database) UpdatePostStatus(postType string, id string, status int8) error {
	tableName := ""
	switch postType {
	case "article":
		tableName = "article_tables"
	case "message":
		tableName = "message_tables"
	case "sharing":
		tableName = "sharing_tables"
	default:
		return errors.New("invalid post type")
	}
	query := d.DB.Table(tableName).Where("id = ?", id)
	if query.Error != nil {
		return query.Error
	}
	if status < 0 || status > 2  {
		return errors.New("status must be 0, 1, or 2")
	}
	if query.RowsAffected == 0 {
		return errors.New("post not found")
	}

	return query.Update("status", status).Error
}
