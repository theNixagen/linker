"use server";

import { db } from "@/src/db/drizzle";
import { users } from "@/src/db/schema";
import { createBucketIfNotExists, s3Client } from "@/src/files/client";
import { eq } from "drizzle-orm";
import { randomUUID } from "node:crypto";

type uploadProfilePicture = {
  profilePicture: any;
  userId: number;
};

export async function uploadPhoto({
  profilePicture,
  userId,
}: uploadProfilePicture) {
  const hashedFileName = randomUUID();
  await createBucketIfNotExists("linker");
  await s3Client.putObject("linker", hashedFileName, profilePicture);
  await db
    .update(users)
    .set({ profilePicture: hashedFileName })
    .where(eq(users.id, userId));
}
