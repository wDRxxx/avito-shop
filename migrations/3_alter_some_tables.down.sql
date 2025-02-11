ALTER TABLE "inventory"
DROP COLUMN "quantity";

ALTER TABLE "users"
ALTER COLUMN "balance" DROP DEFAULT;

ALTER TABLE "users"
DROP CONSTRAINT "users_balance_check";

ALTER TABLE "users"
ADD CONSTRAINT "users_balance_check" CHECK ("balance" > 0);