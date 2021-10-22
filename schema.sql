--6. Проект: Поиск продукта в ближайших аптеках
--7. Клиент заходит в приложение
--Вводит Название продукта(препарата)
--Приложение выводит список карточек продуктов в ближайших аптеках
--Список состоит из Название Города, улицы, название аптеки, название продукта, стоимость продукта, колличество.

--8. Схема

CREATE TABLE stores (
			id            serial PRIMARY KEY,
			name          varchar(200) NOT NULL,
			addr          text NOT NULL,
            street        text NOT NULL,
			name_product  varchar(100),
			
		);

CREATE TABLE products (
			id          serial PRIMARY KEY,
			name        varchar(200) NOT NULL,
			price       integer
		);

CREATE TABLE quantity (
			store_id          integer,
			product_id        integer,
			quantity		  integer,
			FOREIGN KEY (store_id) REFERENCES stores (id),  
			FOREIGN KEY (product_id) REFERENCES products (id)
		);
