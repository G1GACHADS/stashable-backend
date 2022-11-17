BEGIN;
-- view for warehouses list
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

