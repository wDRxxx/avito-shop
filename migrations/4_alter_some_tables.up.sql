ALTER TABLE "users"
ALTER COLUMN "balance" SET DEFAULT 1000;

ALTER TABLE "inventory"
ADD CONSTRAINT "unique_user_item" UNIQUE (user_id, item_id);