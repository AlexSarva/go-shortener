package storagepg

import (
	"AlexSarva/go-shortener/models"
	"AlexSarva/go-shortener/storage"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	database *sqlx.DB
}

func NewPostgresDBConnection(config string) *PostgresDB {
	db, err := sqlx.Connect("postgres", config)
	var schema = `
	CREATE TABLE if not exists public.urls (
		id text,
		short_url text,
		raw_url text primary key,
		user_id text,
		deleted int default 0,
		created timestamp
	);`
	db.MustExec(schema)
	if err != nil {
		log.Fatalln(err)
	}
	return &PostgresDB{
		database: db,
	}
}

func (d *PostgresDB) Ping() bool {
	return d.database.Ping() == nil
}

func (d *PostgresDB) InsertURL(id, rawURL, shortURL, userID string) error {
	URLData := &models.URL{
		ID:       id,
		RawURL:   rawURL,
		ShortURL: shortURL,
		Created:  time.Now(),
		UserID:   userID,
	}

	tx := d.database.MustBegin()
	resInsert, resErr := tx.NamedExec("INSERT INTO public.urls (id, short_url, raw_url, user_id, created) VALUES (:id, :short_url, :raw_url, :user_id, :created) on conflict (raw_url) do nothing ", &URLData)
	affectedRows, _ := resInsert.RowsAffected()
	if affectedRows == 0 {
		return storage.ErrDuplicatePK
	}
	if resErr != nil {
		log.Println(resErr)
	}
	commitErr := tx.Commit()
	if commitErr != nil {
		log.Println(commitErr)
	}
	return nil
}

func (d *PostgresDB) InsertMany(bathURL []models.URL) error {
	_, err := d.database.NamedExec(`INSERT INTO public.urls (id, short_url, raw_url, user_id, created)
        VALUES (:id, :short_url, :raw_url, :user_id, :created) on conflict (raw_url) do nothing`, bathURL)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (d *PostgresDB) GetURL(id string) (*models.URL, error) {
	var getURL models.URL
	err := d.database.Get(&getURL, "SELECT id, short_url, raw_url, user_id, deleted, created FROM public.urls WHERE id=$1", id)
	if err != nil {
		log.Println(err)
	}
	return &getURL, err
}

func (d *PostgresDB) GetURLByRaw(rawURL string) (*models.URL, error) {
	var getURL models.URL
	err := d.database.Get(&getURL, "SELECT id, short_url, raw_url, user_id, deleted, created FROM public.urls WHERE raw_url=$1", rawURL)
	if err != nil {
		log.Println(err)
	}
	return &getURL, err
}

func (d *PostgresDB) GetUserURLs(userID string) ([]models.UserURL, error) {
	var allURLs []models.UserURL
	log.Println(userID)
	err := d.database.Select(&allURLs, "SELECT short_url, raw_url FROM public.urls where user_id=$1", userID)
	if err != nil {
		log.Println(err)
	}
	return allURLs, err
}

func (d *PostgresDB) Delete(userID string, shortURLs []string) error {
	query := "UPDATE public.urls SET deleted=1 WHERE user_id=? AND id IN (?)"
	qry, args, err := sqlx.In(query, userID, shortURLs)
	if err != nil {
		return err
	}

	if _, execErr := d.database.Exec(d.database.Rebind(qry), args...); execErr != nil {
		log.Println(execErr)
		return execErr
	}
	return nil
}
