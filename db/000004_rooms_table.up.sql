BEGIN;
CREATE TABLE IF NOT EXISTS "rooms" (
    "id" bigserial PRIMARY KEY,
    "warehouse_id" bigint REFERENCES "warehouses" ("id") ON DELETE CASCADE,
    "image_url" text NOT NULL,
    "name" varchar NOT NULL,
    "width" decimal NOT NULL,
    "height" decimal NOT NULL,
    "length" decimal NOT NULL,
    "price" decimal NOT NULL
);
ALTER TABLE "rentals"
    ADD COLUMN "room_id" bigint REFERENCES "rooms" ("id") ON DELETE SET NULL ON UPDATE CASCADE;
ALTER TABLE "warehouses"
    ADD COLUMN "rooms_count" integer NOT NULL DEFAULT 0;
--
CREATE OR REPLACE FUNCTION add_warehouse_rooms_count ()
    RETURNS TRIGGER
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE
        warehouses
    SET
        rooms_count = (
            SELECT
                count(*)
            FROM
                rooms
            WHERE
                rooms.warehouse_id = warehouses.id)
    WHERE
        warehouses.id = NEW.warehouse_id;
    RETURN NULL;
END
$$;
--
CREATE TRIGGER add_warehouse_rooms_count
    AFTER INSERT OR UPDATE OR DELETE OR TRUNCATE ON rooms
    FOR EACH STATEMENT
    EXECUTE PROCEDURE add_warehouse_rooms_count ();
-- Re-create view
DROP VIEW IF EXISTS warehouses_list;
CREATE VIEW warehouses_list AS
SELECT
    (
        SELECT
            count(w.id)
        FROM
            warehouses w),
    w.id AS w_id,
    w.address_id,
    w.name,
    w.image_url,
    w.description,
    w.base_price,
    w.email,
    w.phone_number,
    w.created_at,
    w.rooms_count,
    a.id AS a_id,
    a.province,
    a.city,
    a.street_name,
    a.zip_code,
    array_agg(c.name) AS supported_categories
FROM
    warehouses w
    LEFT JOIN addresses AS a ON w.address_id = a.id
    LEFT JOIN warehouse_categories AS wc ON w.id = wc.warehouse_id
    LEFT JOIN categories AS c ON wc.category_id = c.id
GROUP BY
    w.id,
    a.id;
COMMIT;

