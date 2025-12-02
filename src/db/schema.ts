import { integer, pgTable, serial, text } from "drizzle-orm/pg-core";

export const users = pgTable("users", {
  id: serial("id").primaryKey(),
  email: text("email").notNull().unique(),
  name: text("name").notNull(),
  password: text("password").notNull(),
  profilePicture: text("profile_picture"),
  bannerPicture: text("banner_picture"),
});

export const links = pgTable("links", {
  id: serial("id").primaryKey(),
  userId: integer("user_id")
    .references(() => users.id)
    .notNull(),
  url: text("url").notNull(),
  title: text("title").notNull(),
  description: text("description").notNull(),
});
