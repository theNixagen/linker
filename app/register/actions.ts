"use server";

import { db } from "@/src/db/drizzle";
import { users } from "@/src/db/schema";
import * as bcrypt from "bcrypt";

type CreateUser = {
  name: string;
  email: string;
  password: string;
};

export async function createUser({ email, name, password }: CreateUser) {
  const hashed_password = await bcrypt.hash(password, 10);
  await db.insert(users).values({ email, name, password: hashed_password });
}
