CREATE TABLE "users" (
    "id" serial PRIMARY KEY,
    "first_name" varchar NOT NULL,
    "last_name" varchar NOT NULL,
    "email" varchar NOT NULL,
    "password" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);