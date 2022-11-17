BEGIN;
CREATE TABLE IF NOT EXISTS "shipping_types" (
    "id" serial PRIMARY KEY,
    "name" varchar(32) NOT NULL
);
INSERT INTO "shipping_types" ("name")
    VALUES ('pick-up-truck'), ('pick-up-box'), ('van'), ('truck'), ('self-service');
CREATE TABLE IF NOT EXISTS "rental_statuses" (
    "id" serial PRIMARY KEY,
    "name" varchar(32) NOT NULL
);
INSERT INTO "rental_statuses" ("name")
    VALUES ('unpaid'), ('paid'), ('cancelled'), ('returned');
CREATE TABLE IF NOT EXISTS "rentals" (
    "id" bigserial PRIMARY KEY,
    "user_id" bigint REFERENCES "users" ("id") ON DELETE CASCADE,
    "warehouse_id" bigint REFERENCES "warehouses" ("id") ON DELETE CASCADE,
    "category_id" bigint REFERENCES "categories" ("id") ON DELETE CASCADE,
    "shipping_type_id" int NOT NULL REFERENCES "shipping_types" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "rental_status_id" int NOT NULL REFERENCES "rental_statuses" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "image_urls" text[] NOT NULL,
    "name" varchar NOT NULL,
    "description" text NOT NULL,
    "weight" decimal NOT NULL,
    "width" decimal NOT NULL,
    "height" decimal NOT NULL,
    "length" decimal NOT NULL,
    "quantity" int DEFAULT 1 NOT NULL,
    -- if true, then the price is paid annually, otherwise monthly
    "paid_annually" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp(0) without time zone NOT NULL
);
COMMIT;

