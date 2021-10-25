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
            numb_street   INT NOT NULL
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

--CREATE TABLE contacts_user (
            --user_id            INT,
            --contact_id         uuid DEFAULT uuid_generate_v4 (),
            --email              VARCHAR(200) NOT NULL ADD CONSTRAINT valid_email CHECK(VALUE ~ '^\w+@[a-zA-Z_]+?\.[a-zA-Z]{2,3}$'),
            --phone              VARCHAR(100),
            --PRIMARY KEY (contact_id),
            --PRIMARY KEY (user_id) REFERENCES users(id)
--);
