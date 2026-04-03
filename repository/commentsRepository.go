package repository

import (
	"database/sql"
	"errors"
)

func UpdateCommentStatus(db *sql.DB, commentID int64, status string) error {
	res, err := db.Exec(`
		UPDATE comments 
		SET status = $1 
		WHERE id = $2
`, status, commentID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("comment not found")
	}

	return nil
}
