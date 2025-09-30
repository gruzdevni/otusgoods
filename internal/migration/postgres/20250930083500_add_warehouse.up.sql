CREATE TABLE nomenclature(
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(5) NOT NULL
);

INSERT INTO nomenclature (full_name)
VALUES ('Молоток'),
    ('Рубанок'),
    ('Фрезер'),
    ('Станок');
    
CREATE TABLE available_quantity(
    id int PRIMARY KEY,
    available_quantity int NOT NULL
);

INSERT INTO available_quantity (id, available_quantity)
VALUES (1, 1000),
    (2, 500);


CREATE TABLE goods_reservations(
    id SERIAL PRIMARY KEY,
    order_id TEXT NOT NULL,
    nomenclature_id int NOT NULL,
    quantity_reserved int NOT NULL
);

COMMENT ON COLUMN nomenclature.id IS 'id строки';
COMMENT ON COLUMN nomenclature.full_name IS 'наименование товара';
COMMENT ON COLUMN available_quantity.id IS 'id диапазона';
COMMENT ON COLUMN available_quantity.available_quantity IS 'Количество свободных товаров';
COMMENT ON COLUMN goods_reservations.id IS 'id строки';
COMMENT ON COLUMN goods_reservations.order_id IS 'идентификатор заказа';
COMMENT ON COLUMN goods_reservations.nomenclature_id IS 'идентификатор товара';
COMMENT ON COLUMN goods_reservations.quantity_reserved IS 'зарезервированное количество товара';
