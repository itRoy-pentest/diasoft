package storage

import (
	"database/sql"
)

// DiplomaEntry — структура для передачи данных диплома в БД
type DiplomaEntry struct {
	UnivID       string
	StudentID    string
	StudentName  string
	DiplomaNum   string
	IdentityHash string
	Signature    string
	IssueYear    int
}

// SaveDiploma сохраняет запись в Postgres
func SaveDiploma(db *sql.DB, d DiplomaEntry) error {
	query := `INSERT INTO diplomas (
		id, univ_id, student_id, student_name, 
		diploma_number, identity_hash, digital_signature, 
		issue_year, created_at
	) VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, NOW())`

	_, err := db.Exec(query, 
		d.UnivID, d.StudentID, d.StudentName, 
		d.DiplomaNum, d.IdentityHash, d.Signature, 
		d.IssueYear,
	)
	return err
}

// FindByIdentityHash — ТА САМАЯ ФУНКЦИЯ, которой не хватало
func FindByIdentityHash(db *sql.DB, hash string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM diplomas WHERE identity_hash = $1)`
	err := db.QueryRow(query, hash).Scan(&exists)
	return exists, err
}
