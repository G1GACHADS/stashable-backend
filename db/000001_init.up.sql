BEGIN;
CREATE TABLE IF NOT EXISTS "categories" (
    "id" bigserial PRIMARY KEY,
    "name" varchar
);
CREATE TABLE IF NOT EXISTS "addresses" (
    "id" bigserial PRIMARY KEY,
    "province" varchar NOT NULL,
    "city" varchar NOT NULL,
    "street_name" varchar NOT NULL,
    "zip_code" int NOT NULL
);
CREATE TABLE IF NOT EXISTS "users" (
    "id" bigserial PRIMARY KEY,
    "address_id" bigint NOT NULL REFERENCES "addresses" ("id") ON DELETE CASCADE,
    "full_name" varchar NOT NULL,
    "email" varchar NOT NULL,
    "phone_number" varchar NOT NULL,
    "password" text NOT NULL,
    "created_at" timestamp(0) without time zone NOT NULL
);
CREATE UNIQUE INDEX ON "users" ("email", "phone_number");
CREATE TABLE IF NOT EXISTS "warehouses" (
    "id" bigserial PRIMARY KEY,
    "address_id" bigint NOT NULL REFERENCES "addresses" ("id") ON DELETE CASCADE,
    "name" varchar NOT NULL,
    "image_url" text NOT NULL,
    "description" text NOT NULL,
    "base_price" decimal NOT NULL,
    "email" varchar NOT NULL,
    "phone_number" varchar NOT NULL,
    "created_at" timestamp(0) without time zone NOT NULL
);
CREATE TABLE IF NOT EXISTS "warehouse_categories" (
    "id" bigserial PRIMARY KEY,
    "warehouse_id" bigint NOT NULL REFERENCES "warehouses" ("id") ON DELETE CASCADE,
    "category_id" bigint NOT NULL REFERENCES "categories" ("id") ON DELETE CASCADE
);
COMMIT;

