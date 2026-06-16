package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"web-service/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath) // Подключаемся к БД
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Создаем таблицу, если ее еще нет
	stmt, err := db.Prepare(`
    CREATE TABLE IF NOT EXISTS links(
        id INTEGER PRIMARY KEY,
        Short_code TEXT NOT NULL UNIQUE,
        Original_url TEXT NOT NULL,
		Visits INTEGER,
		Created_at TEXT );
    CREATE INDEX IF NOT EXISTS idx_short_code ON links(Short_code);
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

func (s *Storage) SaveLink(Short_code string, OriginalUrl string) (int64, error) {
	const op = "storage.sqlite.SaveLink"

	stmt, err := s.db.Prepare("INSERT INTO links(Original_url, Short_code, Created_at, Visits) VALUES(?, ?, ?, 0)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(OriginalUrl, Short_code, time.Now())
	if err != nil {
		// if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		// 	return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		// }

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetLink(Short_code string, addVisit bool) (models.Link, error) {
	const op = "storage.sqlite.GetLink"

	stmt, err := s.db.Prepare("SELECT Short_code, Original_url, Visits, Created_at FROM links WHERE Short_code = ?")
	if err != nil {
		return models.Link{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var link models.Link
	err = stmt.QueryRow(Short_code).Scan(&link.Short_code, &link.Original_url, &link.Visits, &link.Created_at)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Link{}, models.ErrURLNotFound
		}

		return models.Link{}, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	if addVisit {
		link.Visits++ // Увеличиваем счетчик посещений
		err = s.AddVisit(Short_code, link)
		if err != nil {
			return models.Link{}, fmt.Errorf("%s: add visit: %w", op, err)
		}
	}
	return link, nil
}

func (s *Storage) AddVisit(Short_code string, link models.Link) error {
	const op = "storage.sqlite.AddVisit"

	updateStmt, err := s.db.Prepare("UPDATE links SET Visits = ? WHERE Short_code = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare update statement: %w", op, err)
	}
	_, err = updateStmt.Exec(link.Visits, link.Short_code)
	if err != nil {
		return fmt.Errorf("%s: execute update statement: %w", op, err)
	}
	return nil

}

func (s *Storage) GetLinks(limit int, offset int) ([]models.Link, error) {
	const op = "storage.sqlite.GetLinks"

	rows, err := s.db.Query("SELECT Short_code, Original_url, Visits, Created_at FROM links LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: query: %w", op, err)
	}
	defer rows.Close()

	var links []models.Link
	for rows.Next() {
		var link models.Link
		if err := rows.Scan(&link.Short_code, &link.Original_url, &link.Visits, &link.Created_at); err != nil {
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}
		links = append(links, link)
	}

	return links, nil
}

func (s *Storage) DeleteLink(Short_code string) error {
	const op = "storage.sqlite.DeleteLink"

	stmt, err := s.db.Prepare("DELETE FROM links WHERE Short_code = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	_, err = stmt.Exec(Short_code)

	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return nil
}
