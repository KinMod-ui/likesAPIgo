CREATE TABLE "likes" (
  "user_id" bigint,
  "content_id" bigint,
  "liked" boolean DEFAULT true,
  "update_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("user_id", "content_id")
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "content" (
  "id" bigserial PRIMARY KEY,
  "title" varchar(10) NOT NULL,
  "user_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "notifications" (
  "user_id" bigint PRIMARY KEY
);

CREATE INDEX ON "likes" ("content_id");

ALTER TABLE "likes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "likes" ADD FOREIGN KEY ("content_id") REFERENCES "content" ("id");

ALTER TABLE "content" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "notifications" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
