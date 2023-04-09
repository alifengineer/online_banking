
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "users" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "first_name" varchar(255) NOT NULL,
    "last_name" varchar(255) NOT NULL,
    "phone" varchar NOT NULL UNIQUE,
    "password" varchar(255) NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" INTEGER DEFAULT 0
);

CREATE UNIQUE INDEX "users_phone_deleted_at_unique" ON "users" ("phone", "deleted_at");

CREATE TABLE IF NOT EXISTS "accounts" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "balance" numeric NOT NULL DEFAULT 0.0,
    "user_id" UUID NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" INTEGER DEFAULT 0

    CONSTRAINT "positive_balance"
        CHECK ("balance" >= 0.0),

    CONSTRAINT "accounts_user_id_fkey"
        FOREIGN KEY ("user_id")
        REFERENCES "users" ("guid")
);

CREATE UNIQUE INDEX "accounts_user_id_deleted_at_unique" ON "accounts" ("user_id", "deleted_at");

CREATE TABLE IF NOT EXISTS "transactions" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "account_id" UUID NOT NULL,
    "recipient_id" UUID NOT NULL,
    "transaction_type" varchar(255) NOT NULL,
    "transaction_amount" numeric NOT NULL,
    "approved" BOOLEAN NOT NULL DEFAULT false,
    "done" BOOLEAN NOT NULL DEFAULT false,
    "done_timestamp" TIMESTAMP WITH TIME ZONE,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL

    CONSTRAINT "positive_transaction"
        CHECK ("transaction_amount" > 0.0),

    CONSTRAINT "transactions_account_id_fkey"
        FOREIGN KEY ("account_id")
        REFERENCES "accounts" ("guid"),
    
    CONSTRAINT "transactions_recipient_id_fkey"
        FOREIGN KEY ("recipient_id")
        REFERENCES "accounts" ("guid")
);