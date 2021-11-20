--6. Проект: Поиск продукта в ближайших аптеках
--7. Клиент заходит в приложение
--Вводит Название продукта(препарата)
--Приложение выводит список карточек продуктов в ближайших аптеках
--Список состоит из Название Города, улицы, название аптеки, название продукта, стоимость продукта, колличество.

--8. Схема

CREATE TABLE stores (
			id            INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			name          VARCHAR(200) NOT NULL,
			city          VARCHAR(200) NOT NULL,
            street        VARCHAR(200) NOT NULL,
            numb_street   VARCHAR(200) NOT NULL,
            email         VARCHAR(200) NOT NULL,
            phone         VARCHAR(100),
            lon           VARCHAR(100),
			lat           VARCHAR(100)
		);

CREATE TABLE products (
			id          INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			name        VARCHAR(200) NOT NULL,
			price       MONEY CONSTRAINT positive_price CHECK (price::numeric > 0)
		);

CREATE TABLE quantity (
			store_id          INT,
			product_id        INT,
			quantity		  INT,
			FOREIGN KEY (store_id) REFERENCES stores (id),  
			FOREIGN KEY (product_id) REFERENCES products (id)
		);

CREATE TABLE users(
            id            INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
            first_name    VARCHAR(200) NOT NULL,
            last_name     VARCHAR(200) NOT NULL,
            city_user          VARCHAR(200) NOT NULL,
            street_user        VARCHAR(200) NOT NULL,
            numb_street_user   INT NOT NULL
        );


--CREATE DOMAIN domain_email AS citext
--CHECK(
-- VALUE ~ '^\w+@[a-zA-Z_]+?\.[a-zA-Z]{2,3}$'
--);

CREATE TABLE contacts_user (
            user_id            INT,
            contact_id         uuid DEFAULT uuid_generate_v4 (),
            email              VARCHAR(200) NOT NULL,
            phone              VARCHAR(100),
            PRIMARY KEY (contact_id),
            FOREIGN KEY (user_id) REFERENCES users(id)
);

create table analog_prod ( --для тестирования
	name VARCHAR(50),
	city VARCHAR(50),
	email VARCHAR(50),
	phone VARCHAR(50),
	price VARCHAR(50)
);

create index if not exists idx_store_id_product_id on quantity (store_id, product_id);
CREATE INDEX concurrently analog_email_idx ON analog_prod USING btree (email text_pattern_ops) include (phone);
CREATE INDEX concurrently analog_name_idx ON analog_prod (name);
