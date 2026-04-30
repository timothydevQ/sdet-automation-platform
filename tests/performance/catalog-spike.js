import http from "k6/http";
import { sleep } from "k6";

const API = __ENV.API_BASE || "http://localhost:8080";

export const options = {
  scenarios: {
    spike: {
      executor: "ramping-arrival-rate",
      startRate: 5,
      timeUnit: "1s",
      preAllocatedVUs: 50,
      maxVUs: 100,
      stages: [
        { duration: "10s", target: 5 },
        { duration: "10s", target: 200 },
        { duration: "20s", target: 200 },
        { duration: "10s", target: 0 },
      ],
    },
  },
  thresholds: {
    http_req_duration: ["p(95)<500"],
    http_req_failed: ["rate<0.01"],
  },
};

export default function () {
  http.get(`${API}/catalog/products`);
  sleep(0.1);
}
