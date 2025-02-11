ALTER TABLE "inventory"
DROP COLUMN "quantity";

ALTER TABLE "users"
ALTER COLUMN "balance" DROP DEFAULT;

ALTER TABLE "users"
DROP CONSTRAINT "users_balance_check";

UPDATE "users"
SET "balance" = 1
WHERE "balance" = 0;

ALTER TABLE "users"
ADD CONSTRAINT "users_balance_check" CHECK ("balance" > 0);