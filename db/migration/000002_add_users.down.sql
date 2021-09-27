ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS  "citizen_currency_key";

ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS  "accounts_owner_fkey";

ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS  "accounts_citizenship_fkey";

DROP  TABLE IF EXISTS "users";