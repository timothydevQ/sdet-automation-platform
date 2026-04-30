import http from "k6/http";
import { check, sleep } from "k6";
import { Trend, Rate } from "k6/metrics";

const checkoutLatency = new Trend("checkout_latency_ms");
const checkoutErrors = new Rate("checkout_errors");

const API = __ENV.API_BASE || "http://localhost:8080";

export const options = {
  scenarios: {
    ramp: {
      executor: "ramping-vus",
      startVUs: 1,
      stages: [
        { duration: "30s", target: 10 },
        { duration: "1m", target: 25 },
        { duration: "30s", target: 0 },
      ],
    },
  },
  thresholds: {
    http_req_failed: ["rate<0.02"],
    checkout_latency_ms: ["p(95)<800", "p(99)<1500"],
    checkout_errors: ["rate<0.05"],
  },
};

function rand(prefix) {
  return prefix + Math.random().toString(36).slice(2, 12);
}

function register() {
  const email = rand("u") + "@test.local";
  const r = http.post(`${API}/auth/register`,
    JSON.stringify({ email, password: "Hunter22!" }),
    { headers: { "Content-Type": "application/json" } });
  return r.json("token");
}

export default function () {
  const token = register();
  const headers = { "Content-Type": "application/json", Authorization: `Bearer ${token}` };

  http.post(`${API}/cart/items`, JSON.stringify({ sku: "SKU-001", qty: 1 }), { headers });

  const start = Date.now();
  const idem = rand("idem");
  const r = http.post(`${API}/checkout`,
    JSON.stringify({ coupon: "", card_token: "tok_test_visa" }),
    { headers: { ...headers, "Idempotency-Key": idem } });
  checkoutLatency.add(Date.now() - start);

  const ok = check(r, {
    "checkout 201": (res) => res.status === 201,
  });
  checkoutErrors.add(!ok);

  sleep(0.5);
}
