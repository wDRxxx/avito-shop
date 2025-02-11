ALTER TABLE "inventory"
ADD COLUMN "quantity" INTEGER NOT NULL default 1;

ALTER TABLE "users"
ALTER COLUMN "balance" SET DEFAULT 0;

ALTER TABLE "users"
DROP CONSTRAINT "users_balance_check";

ALTER TABLE "users"
ADD CONSTRAINT "users_balance_check" CHECK ("balance" >= 0);