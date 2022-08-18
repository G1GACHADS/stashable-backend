BEGIN;
DROP TRIGGER IF EXISTS refresh_warehouses_list_mat_view ON categories;
DROP TRIGGER IF EXISTS refresh_warehouses_list_mat_view ON warehouses;
DROP FUNCTION IF EXISTS refresh_warehouses_list_mat_view;
DROP MATERIALIZED VIEW IF EXISTS warehouses_list CASCADE;
COMMIT;

