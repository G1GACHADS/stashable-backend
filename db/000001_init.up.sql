BEGIN;
CREATE TABLE "users" (
    "id" BIGSERIAL PRIMARY KEY,
    "address_id" bigint NOT NULL,
    "full_name" varchar NOT NULL,
    "email" varchar NOT NULL,
    "phone_number" varchar NOT NULL,
    "password" text NOT NULL,
    "created_at" timestamp(0) without time zone NOT NULL
);
CREATE UNIQUE INDEX ON "users" ("email", "phone_number");
CREATE TABLE "warehouses" (
    "id" BIGSERIAL PRIMARY KEY,
    "address_id" bigint NOT NULL,
    "name" varchar NOT NULL,
    "image_url" text NOT NULL,
    "description" text NOT NULL,
    "base_price" decimal NOT NULL,
    "email" varchar NOT NULL,
    "phone_number" varchar NOT NULL,
    "created_at" timestamp(0) without time zone NOT NULL
);
CREATE TABLE "warehouse_categories" (
    "id" BIGSERIAL PRIMARY KEY,
    "warehouse_id" bigint NOT NULL,
    "category_id" bigint NOT NULL
);
CREATE TABLE "items" (
    "id" BIGSERIAL PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "warehouse_id" bigint NOT NULL,
    "category_id" bigint NOT NULL,
    "name" varchar NOT NULL,
    "weight" decimal NOT NULL,
    "width" decimal NOT NULL,
    "height" decimal NOT NULL,
    "length" decimal NOT NULL,
    "quantity" int DEFAULT 1 NOT NULL,
    "created_at" timestamp(0) without time zone NOT NULL
);
CREATE TABLE "categories" ("id" BIGSERIAL PRIMARY KEY, "name" varchar);
CREATE TABLE "addresses" (
    "id" BIGSERIAL PRIMARY KEY,
    "province" varchar NOT NULL,
    "city" varchar NOT NULL,
    "street_name" varchar NOT NULL,
    "zip_code" int NOT NULL
);
COMMENT ON COLUMN "items"."weight" IS 'Use kg unit for now';
COMMENT ON COLUMN "items"."width" IS 'Use cm unit for now';
COMMENT ON COLUMN "items"."height" IS 'Use cm unit for now';
ALTER TABLE "users"
ADD FOREIGN KEY ("address_id") REFERENCES "addresses" ("id") ON DELETE CASCADE;
ALTER TABLE "warehouses"
ADD FOREIGN KEY ("address_id") REFERENCES "addresses" ("id") ON DELETE CASCADE;
ALTER TABLE "warehouse_categories"
ADD FOREIGN KEY ("warehouse_id") REFERENCES "warehouses" ("id") ON DELETE CASCADE;
ALTER TABLE "warehouse_categories"
ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON DELETE CASCADE;
ALTER TABLE "items"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "items"
ADD FOREIGN KEY ("warehouse_id") REFERENCES "warehouses" ("id") ON DELETE CASCADE;
ALTER TABLE "items"
ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON DELETE CASCADE;
COMMIT;