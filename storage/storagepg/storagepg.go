package storagepg

import (
	"AlexSarva/go-shortener/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type PostgresDB struct {
	database *sqlx.DB
}

func NewPostgresDBConnection(config string) *PostgresDB {
	//"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	db, err := sqlx.Connect("postgres", config)
	var schema = `
	CREATE TABLE if not exists public.urls (
		id text primary key,
		short_url text,
		raw_url text,
		user_id text,
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

func (d *PostgresDB) InsertURL(id, rawURL, baseURL, userID string) error {
	URLData := &models.URL{
		ID:       id,
		RawURL:   rawURL,
		ShortURL: baseURL + "/" + id,
		Created:  time.Now(),
		UserID:   userID,
	}

	tx := d.database.MustBegin()
	_, resErr := tx.NamedExec("INSERT INTO urls (id, short_url, raw_url, user_id, created) VALUES (:id, :short_url, :raw_url, :user_id, :created)", &URLData)
	if resErr != nil {
		log.Println(resErr)
	}
	commitErr := tx.Commit()
	if commitErr != nil {
		log.Println(commitErr)
	}
	return nil
}

func (d *PostgresDB) GetURL(id string) (*models.URL, error) {
	var getURL models.URL
	err := d.database.Get(&getURL, "SELECT * FROM public.urls WHERE id=$1", id)
	if err != nil {
		log.Println(err)
	}
	return &getURL, err
}

func (d *PostgresDB) GetUserURLs(userID string) ([]models.UserURL, error) {
	var allURLs []models.UserURL
	log.Println(userID)
	err := d.database.Select(&allURLs, "SELECT short_url, raw_url FROM public.urls where user_id=$1", userID)
	return allURLs, err
}
