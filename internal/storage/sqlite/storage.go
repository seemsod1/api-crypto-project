package sqlite

import (
	"api-crypto-project/internal/http-server/handlers/mail"
	"api-crypto-project/internal/http-server/handlers/mail/sendEmails"
	"api-crypto-project/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Storage struct {
	db *sql.DB
}

func NewDB(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS subscribers(
		id INTEGER PRIMARY KEY,
		mail TEXT NOT NULL UNIQUE);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveMail(mailToSave string) (int64, error) {
	const op = "storage.mysql.SaveMail"

	stmt, err := s.db.Prepare("INSERT INTO subscribers(mail) VALUES (?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(mailToSave)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.MailExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}
	return id, nil
}

func (s *Storage) SendEmails(crypto mail.Credits) error {
	const op = "storage.mysql.SendEmails"

	rows, err := s.db.Query("SELECT mail from subscribers where id = ?", 1)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)
	var mailToSend string
	for rows.Next() {
		if err = rows.Scan(&mailToSend); err != nil {
			log.Print("Failed to read mail")
		} else {
			if err := sendEmails.SendToSingleMail(mailToSend, crypto); err != nil {
				return err
			}
			log.Print("Successfully sent")
		}
	}
	return nil
}
