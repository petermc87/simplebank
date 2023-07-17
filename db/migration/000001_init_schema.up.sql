CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'Can be a negative or positive value.';

COMMENT ON COLUMN "transfers"."amount" IS 'Must be a positive number value.';


-- CREATE TABLE "accounts" (
--   "id" bigserial PRIMARY KEY,
--   "owner" varchar NOT NULL,
--   "balance" bigint NOT NULL,
--   "currency" varchar NOT NULL,
--   "created_at" timestamptz NOT NULL DEFAULT (now())
-- );

-- CREATE TABLE "entries" (
--     "id" bigserial PRIMARY KEY,
--     "account_id" bigint NOT NULL,
--     "amount" bigint NOT NULL,
--     "created_at" timestamptz NOT NULL DEFAULT (now())
-- )

-- CREATE TABLE "transfers" (
--     "id" bigserial PRIMARY KEY,
--     "from_account_id" bigint NOT NULL,
--     "to_account_id" bigint NOT NULL,
--     "amount" bigint NOT NULL,
--     "created_at" timestamptz NOT NULL DEFAULT (now())
-- )

-- -- Making the entries "account_id" equal to the "id" in the accounts table.
-- ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");
-- ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");
-- ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

 
-- -- Create indexes for the following so that they can be found easier.
-- CREATE INDEX ON "accounts" ("owner");
-- CREATE INDEX ON "entries" ("account_id");
-- CREATE INDEX ON "transfers" ("from_account_id");
-- CREATE INDEX ON "transfers" ("to_account_id");
-- CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

-- -- Comments on a specifed column.
-- COMMENT ON COLUMN "entries"."amount" IS "Can be negative or positive.";
-- COMMENT ON COLUMN "transfers"."amount" IS "Must be positive";