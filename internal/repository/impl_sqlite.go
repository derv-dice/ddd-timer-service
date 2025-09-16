package repository

import (
	"context"
	"database/sql"
	"ddd-timer-service/models"
	"errors"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

type implSQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(path string, createNew bool) (Repository, error) {
	if path == "" {
		return nil, errors.New("path can't be empty")
	}

	// Проверяем, что указанный файл БД существует
	var shouldMigrate bool

	if _, err := os.ReadFile(path); err != nil {
		if errors.Is(err, os.ErrNotExist) && createNew {
			_, err = os.Create(path)
			if err != nil {
				return nil, err
			}

			shouldMigrate = true
		}
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if shouldMigrate {
		if _, err = db.Exec(sqliteCreateTables); err != nil {
			return nil, err
		}
	}

	return &implSQLiteRepository{db: db}, nil
}

func (r *implSQLiteRepository) SaveUser(ctx context.Context, user *models.User) error {
	serveFromUnix, serveToUnix, birthDate := int64(0), int64(0), int64(0)

	if !user.ServeFrom.IsZero() {
		serveFromUnix = user.ServeFrom.Unix()
	}

	if !user.ServeTo.IsZero() {
		serveToUnix = user.ServeTo.Unix()
	}

	if !user.BirthDate.IsZero() {
		birthDate = user.BirthDate.Unix()
	}

	_, err := r.db.ExecContext(ctx, `
	insert into users (id, date_from, date_to, birth_date) 
	values ($1, $2, $3, $4) 
	on conflict (id) do 
	update set date_from=$2, date_to=$3, birth_date=$4`,
		user.ID, serveFromUnix, serveToUnix, birthDate)

	return err
}

func (r *implSQLiteRepository) LoadUser(ctx context.Context, userID int64) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, `
	select id, date_from, date_to, birth_date
	from users
	where id = $1`, userID)

	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, row.Err()
	}

	serveFromUnix, serveToUnix, birthDate := int64(0), int64(0), int64(0)

	u := new(models.User)
	err := row.Scan(&u.ID, &serveFromUnix, &serveToUnix, &birthDate)
	if err != nil {
		return nil, err
	}

	u.ServeFrom = time.Unix(serveFromUnix, 0)
	u.ServeTo = time.Unix(serveToUnix, 0)
	u.BirthDate = time.Unix(birthDate, 0)

	return u, nil
}

func (r *implSQLiteRepository) DeleteUser(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, `	delete from users where id=$1`, userID)
	return err
}

func (r *implSQLiteRepository) LoadAllUsers(ctx context.Context) ([]*models.User, error) {
	rows, err := r.db.QueryContext(ctx, `
	select id, date_from, date_to, birth_date
	from users`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		u := new(models.User)
		serveFromUnix, serveToUnix, birthDate := int64(0), int64(0), int64(0)

		err = rows.Scan(&u.ID, &serveFromUnix, &serveToUnix, &birthDate)
		if err != nil {
			return nil, err
		}

		u.ServeFrom = time.Unix(serveFromUnix, 0)
		u.ServeTo = time.Unix(serveToUnix, 0)
		u.BirthDate = time.Unix(birthDate, 0)

		users = append(users, u)
	}

	return users, nil
}

const sqliteCreateTables = `
create table users
(
    id         INTEGER not null
        constraint users_pk
            primary key
            on conflict fail,
    date_from  INTEGER not null,
    date_to    INTEGER not null,
    birth_date INTEGER
);`
