package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAddUpdatedAt, downAddUpdatedAt)
}

func upAddUpdatedAt(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS stores (
			id            INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			name          VARCHAR(200) NOT NULL,
			city          VARCHAR(200) NOT NULL,
            street        VARCHAR(200) NOT NULL,
            numb_street   VARCHAR(200) NOT NULL,
            email         VARCHAR(200) NOT NULL,
            phone         VARCHAR(100),
            lon           VARCHAR(100),
			lat           VARCHAR(100)
		);`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id          INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			name        VARCHAR(200) NOT NULL,
			price       MONEY CONSTRAINT positive_price CHECK (price::numeric > 0)
		);`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS quantity (
			store_id          INT,
			product_id        INT,
			quantity		  INT,
			FOREIGN KEY (store_id) REFERENCES stores (id),  
			FOREIGN KEY (product_id) REFERENCES products (id)
		);`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id                 INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
            first_name         VARCHAR(200) NOT NULL,
            last_name          VARCHAR(200) NOT NULL,
            city_user          VARCHAR(200) NOT NULL,
            street_user        VARCHAR(200) NOT NULL,
            numb_street_user   INT NOT NULL
		);`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS contacts_user (
			user_id            INT,
            contact_id         uuid DEFAULT uuid_generate_v4 (),
            email              VARCHAR(200) NOT NULL,
            phone              VARCHAR(100),
            PRIMARY KEY (contact_id),
            FOREIGN KEY (user_id) REFERENCES users(id)
		);`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
			CREATE INDEX IF NOT EXISTS idx_store_id_product_id ON quantity (store_id, product_id);
		`)
	if err != nil {
		return err
	}


	return nil
}

func downAddUpdatedAt(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec("drop table stores;")
	if err != nil {
		return err
	}
	_, err = tx.Exec("drop table products;")
	if err != nil {
		return err
	}
	_, err = tx.Exec("drop table quantity;")
	if err != nil {
		return err
	}
	_, err = tx.Exec("drop table users;")
	if err != nil {
		return err
	}
	_, err = tx.Exec("drop table contacts_user;")
	if err != nil {
		return err
	}
	return nil
}
