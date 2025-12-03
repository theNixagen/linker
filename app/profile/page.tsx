"use client";
import { redirect } from "next/navigation";
import { useSession } from "next-auth/react";

export default function Profile() {
  const { data: sessionData } = useSession();

  if (!sessionData || !sessionData.user) {
    redirect("/login");
  }

  return (
    <div>
      <h1>Profile</h1>
      <p>Welcome, {}!</p>
    </div>
  );
}
