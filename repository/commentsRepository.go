package repository

import (
	"comment-service/models"
	"database/sql"
	"errors"
)

func CreateComment(db *sql.DB, reviewID, userID int64, req models.CreateCommentRequest) error {
	_, err := db.Exec(`
		INSERT INTO comments (review_id, user_id, comment_text, status, created_at, prev_comment_id)
		VALUES ($1, $2, $3, $4, NOW(), $5)
`, reviewID, userID, req.Text, "moderating", SafeDeref(req.PrevCommentID))
	if err != nil {
		return err
	}

	return nil
}

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

func UpdateComment(db *sql.DB, commentID int64, req models.UpdateCommentRequest) error {
	res, err := db.Exec(`
		UPDATE comments
		SET comment_text = $1, status = 'moderating', 
		    edited_flag = true, edited_at = NOW()
		WHERE id = $2
`, req.Text, commentID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("comment not found")
	}

	return nil
}

func GetUserIdAndReviewIDByComment(db *sql.DB, commentID int64) (int64, int64, error) {
	var userID, reviewID int64
	err := db.QueryRow(`SELECT user_id, review_id FROM comments WHERE id = $1`, commentID).Scan(&userID, &reviewID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, errors.New("incorrect commentID")
		}
		return 0, 0, err
	}

	return userID, reviewID, nil
}

func DeleteComment(db *sql.DB, commentID int64) error {
	res, err := db.Exec(`
		DELETE FROM comments
		WHERE id = $1
`, commentID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("comment not found")
	}

	return nil
}

func DeleteCommentsByReview(db *sql.DB, reviewID int64) error {
	_, err := db.Exec(`
		DELETE FROM comments
		WHERE review_id = $1
`, reviewID)
	if err != nil {
		return err
	}

	return nil
}

func GetCommentsByReview(db *sql.DB, reviewID int64) ([]models.Comment, error) {
	rows, err := db.Query(`
		SELECT id, review_id, user_id, comment_text, status, created_at, 
		       COALESCE(prev_comment_id, 0), edited_flag, COALESCE(edited_at, created_at)
		FROM comments
		WHERE review_id = $1
`, reviewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.ReviewID,
			&comment.UserID,
			&comment.Text,
			&comment.Status,
			&comment.CreatedAt,
			&comment.PrevCommentID,
			&comment.EditedFlag,
			&comment.EditedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func SafeDeref[T any](v *T) any {
	if v == nil {
		return nil
	}
	return *v
}
