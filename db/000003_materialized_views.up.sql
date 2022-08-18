BEGIN;
-- Materialized view for warehouses list
CREATE MATERIALIZED VIEW warehouses_list AS
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
    array_agg(c.name)
FROM
    warehouses w
    LEFT JOIN addresses AS a ON w.address_id = a.id
    LEFT JOIN warehouse_categories AS wc ON w.id = wc.warehouse_id
    LEFT JOIN categories AS c ON wc.category_id = c.id
GROUP BY
    w.id,
    a.id;
-- Refresh function
CREATE OR REPLACE FUNCTION refresh_warehouses_list_mat_view ()
    RETURNS TRIGGER
    LANGUAGE plpgsql
    AS $$
BEGIN
    REFRESH MATERIALIZED VIEW warehouses_list;
    RETURN NULL;
END
$$;
-- Triggers
CREATE TRIGGER refresh_warehouses_list_mat_view
    AFTER INSERT OR UPDATE OR DELETE OR TRUNCATE ON warehouses
    FOR EACH STATEMENT
    EXECUTE PROCEDURE refresh_warehouses_list_mat_view ();
CREATE TRIGGER refresh_warehouses_list_mat_view
    AFTER INSERT OR UPDATE OR DELETE OR TRUNCATE ON categories
    FOR EACH STATEMENT
    EXECUTE PROCEDURE refresh_warehouses_list_mat_view ();
--
COMMIT;

