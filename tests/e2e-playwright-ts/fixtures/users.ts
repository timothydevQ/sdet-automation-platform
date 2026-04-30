import { APIRequestContext } from "@playwright/test";
import { randomUUID } from "crypto";

const API = process.env.API_BASE || "http://localhost:8080";

export async function createCustomer(req: APIRequestContext) {
  const email = `u${randomUUID().replace(/-/g, "").slice(0, 10)}@test.local`;
  const password = "Hunter22!";
  const r = await req.post(`${API}/auth/register`, { data: { email, password } });
  return { email, password, ...(await r.json()) };
}

export async function createAdmin(req: APIRequestContext) {
  const email = `a${randomUUID().replace(/-/g, "").slice(0, 10)}@admin.local`;
  const password = "Hunter22!";
  const r = await req.post(`${API}/auth/register`, { data: { email, password } });
  return { email, password, ...(await r.json()) };
}
