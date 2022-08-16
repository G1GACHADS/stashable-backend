BEGIN;
CREATE TABLE IF NOT EXISTS "rental_types" ("name" varchar(18) NOT NULL PRIMARY KEY);
INSERT INTO "rental_types" ("name")
VALUES ('self-storage'),
    ('disposal');
CREATE TABLE IF NOT EXISTS "rental_statuses" ("name" varchar(36) NOT NULL PRIMARY KEY);
INSERT INTO "rental_statuses" ("name")
VALUES ('unpaid'),
    ('paid'),
    ('cancelled'),
    ('returned');
CREATE TABLE IF NOT EXISTS "rentals" (
    "id" BIGSERIAL PRIMARY KEY,
    "user_id" BIGINT REFERENCES "users" ("id") ON DELETE CASCADE,
    "warehouse_id" BIGINT REFERENCES "warehouses" ("id") ON DELETE CASCADE,
    "category_id" BIGINT REFERENCES "categories" ("id") ON DELETE CASCADE,
    "image_urls" text[] NOT NULL,
    "name" varchar NOT NULL,
    "description" text NOT NULL,
    "weight" decimal NOT NULL,
    "width" decimal NOT NULL,
    "height" decimal NOT NULL,
    "length" decimal NOT NULL,
    "quantity" int DEFAULT 1 NOT NULL,
    -- if true, then the price is paid annually, otherwise monthly
    "paid_annually" boolean NOT NULL DEFAULT false,
    "type" varchar(18) NOT NULL REFERENCES "rental_types" ("name") ON DELETE CASCADE ON UPDATE CASCADE DEFAULT 'self-storage',
    "status" varchar(36) NOT NULL REFERENCES "rental_statuses" ("name") ON DELETE CASCADE ON UPDATE CASCADE DEFAULT 'pending',
    "created_at" timestamp(0) without time zone NOT NULL
);
COMMIT;