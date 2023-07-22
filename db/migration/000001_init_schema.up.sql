CREATE TABLE "users" (
    "id" serial PRIMARY KEY,
    "firstname" varchar NOT NULL,
    "full_name" varchar NOT NULL,
    "email" varchar NOT NULL,
    "password" varchar,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);