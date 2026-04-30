const BASE = (import.meta as any).env.VITE_API_BASE || "http://localhost:8080";

function token() {
  return localStorage.getItem("token");
}

export async function api<T>(path: string, opts: RequestInit = {}): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(opts.headers as Record<string, string>),
  };
  const t = token();
  if (t) headers.Authorization = `Bearer ${t}`;
  const res = await fetch(`${BASE}${path}`, { ...opts, headers });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || res.statusText);
  }
  return res.status === 204 ? (undefined as T) : res.json();
}

export const auth = {
  login: (email: string, password: string) =>
    api<{ token: string; role: string }>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    }),
  register: (email: string, password: string) =>
    api<{ token: string; role: string }>("/auth/register", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    }),
};

export const catalog = {
  list: () => api<any[]>("/catalog/products"),
};

export const cart = {
  get: () => api<any[]>("/cart"),
  add: (sku: string, qty: number) =>
    api("/cart/items", { method: "POST", body: JSON.stringify({ sku, qty }) }),
  remove: (sku: string) => api(`/cart/items/${sku}`, { method: "DELETE" }),
};

export const orders = {
  checkout: (coupon: string, cardToken: string, idemKey?: string) =>
    api<any>("/checkout", {
      method: "POST",
      headers: idemKey ? { "Idempotency-Key": idemKey } : {},
      body: JSON.stringify({ coupon, card_token: cardToken }),
    }),
  get: (id: number) => api<any>(`/orders/${id}`),
  adminList: () => api<any[]>("/admin/orders"),
  adminRefund: (id: number) =>
    api(`/admin/orders/${id}/refund`, { method: "POST" }),
};
