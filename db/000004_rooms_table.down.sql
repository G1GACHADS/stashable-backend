BEGIN;
DROP MATERIALIZED VIEW IF EXISTS warehouses_list;
DROP TRIGGER IF EXISTS add_warehouse_rooms_count ON "rooms";
DROP FUNCTION IF EXISTS add_warehouse_rooms_count;
DROP TABLE IF EXISTS "rooms" CASCADE;
ALTER TABLE "rentals"
    DROP COLUMN "room_id";
ALTER TABLE "warehouses"
    DROP COLUMN "rooms_count";
COMMIT;

