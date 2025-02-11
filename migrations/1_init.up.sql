CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,
    "username" VARCHAR(100) UNIQUE NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "balance" INTEGER NOT NULL CHECK ("balance" > 0),

    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_username ON "users"("username");

CREATE TABLE IF NOT EXISTS "items" (
    "id" SERIAL PRIMARY KEY,
    "title" VARCHAR(100) UNIQUE NOT NULL,
    "price" INTEGER NOT NULL,

    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_items_title ON "items"("title");

CREATE TABLE IF NOT EXISTS "inventory" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER,
    "item_id" INTEGER,

    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT "fk_user"
        FOREIGN KEY ("user_id")
            REFERENCES "users"("id")
            ON DELETE CASCADE,

    CONSTRAINT "fk_item"
        FOREIGN KEY ("item_id")
            REFERENCES "items"("id")
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "transactions" (
    "id" SERIAL PRIMARY KEY,
    "type" BOOLEAN NOT NULL,
    "sender" INTEGER,
    "recipient" INTEGER,
    "amount" INTEGER NOT NULL,

    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT "fk_sender"
        FOREIGN KEY ("sender")
            REFERENCES "users"("id")
            ON DELETE CASCADE,

    CONSTRAINT "fk_recipient"
        FOREIGN KEY ("recipient")
            REFERENCES "users"("id")
            ON DELETE CASCADE
);