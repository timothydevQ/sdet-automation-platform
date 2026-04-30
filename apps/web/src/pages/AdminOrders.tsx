import React, { useEffect, useState } from "react";
import { orders } from "../lib/api";

export default function AdminOrders() {
  const [list, setList] = useState<any[]>([]);
  const [err, setErr] = useState("");

  useEffect(() => {
    orders.adminList().then(setList).catch((e) => setErr(e.message));
  }, []);

  if (err) return <p data-testid="admin-error" role="alert">{err}</p>;

  return (
    <section data-testid="admin-orders">
      <h1>Orders</h1>
      <table>
        <thead>
          <tr><th>ID</th><th>User</th><th>Status</th><th>Total</th><th></th></tr>
        </thead>
        <tbody>
          {list.map((o) => (
            <tr key={o.id} data-testid={`order-row-${o.id}`}>
              <td>{o.id}</td>
              <td>{o.user_id}</td>
              <td>{o.status}</td>
              <td>${(o.total_cents / 100).toFixed(2)}</td>
              <td>
                <button
                  data-testid={`refund-${o.id}`}
                  onClick={() => orders.adminRefund(o.id).then(() => location.reload())}
                >
                  Refund
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
}
