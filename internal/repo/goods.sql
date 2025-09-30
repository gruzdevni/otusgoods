-- name: GetAvailableGoods :many
SELECT id, available_quantity FROM available_quantity WHERE id = ANY(@id::varchar[]);

-- name: ReserveGoodsForOrder :exec
INSERT INTO goods_reservations (order_id, nomenclature_id, quantity_reserved) VALUES ($1, $2, $3);

-- name: UnreserveGoodsForOrder :exec
DELETE FROM goods_reservations WHERE order_id = @order_id;

-- name: DecreaseAvailableGoods :exec
UPDATE available_quantity SET available_quantity = available_quantity - @request_quantity WHERE id = @id;

-- name: IncreaseAvailableGoods :exec
UPDATE available_quantity SET available_quantity = available_quantity + @request_quantity WHERE id = @id;

-- name: CheckOrderReserve :many
SELECT * FROM goods_reservations WHERE order_id = @order_id;
