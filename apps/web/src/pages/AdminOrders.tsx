import React, { useEffect, useState } from "react";
import { orders } from "../lib/api";

export default function AdminOrders() {
  const [list, setList] = useState<any[]>([]);
  const [err, setErr] = useState("");

  function load() {
    orders.adminList().then(setList).catch((e) => setErr(e.message));
  }

  useEffect(() => { load(); }, []);

  if (err) return <p data-testid="admin-error" role="alert">{err}</p>;

  return (
    <section data-testid="admin-orders">
      <h1>Orders <span style={{ fontSize: "1rem", color: "var(--gray-400)", fontWeight: 400 }}>({list.length})</span></h1>
      <table>
        <thead>
          <tr>
            <th>#</th>
            <th>User</th>
            <th>Status</th>
            <th>Total</th>
            <th>Date</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {list.map((o) => (
            <tr key={o.id} data-testid={`order-row-${o.id}`}>
              <td style={{ color: "var(--gray-400)", fontFamily: "monospace" }}>{o.id}</td>
              <td>{o.user_id}</td>
              <td>
                <span className={`status-badge status-${o.status}`}>{o.status}</span>
              </td>
              <td style={{ fontWeight: 600 }}>${(o.total_cents / 100).toFixed(2)}</td>
              <td style={{ color: "var(--gray-400)", fontSize: "0.82rem" }}>
                {new Date(o.created_at).toLocaleDateString()}
              </td>
              <td>
                {o.status === "paid" && (
                  <button
                    data-testid={`refund-${o.id}`}
                    onClick={() => orders.adminRefund(o.id).then(load)}
                  >
                    Refund
                  </button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      {list.length === 0 && (
        <p style={{ textAlign: "center", color: "var(--gray-400)", padding: "3rem" }}>No orders yet.</p>
      )}
    </section>
  );
}
