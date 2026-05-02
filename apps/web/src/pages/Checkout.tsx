import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { orders } from "../lib/api";

export default function Checkout() {
  const [coupon, setCoupon] = useState("");
  const [card, setCard] = useState("tok_test_visa");
  const [err, setErr] = useState("");
  const [loading, setLoading] = useState(false);
  const nav = useNavigate();

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setErr("");
    setLoading(true);
    try {
      const o = await orders.checkout(coupon, card, crypto.randomUUID());
      nav(`/`);
      alert(`Order #${o.id} placed! Total: $${(o.total_cents / 100).toFixed(2)}`);
    } catch (e: any) {
      setErr(e.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <section data-testid="checkout-page">
      <div className="form-card">
        <h1 style={{ marginBottom: "1.5rem", textAlign: "center" }}>Checkout</h1>
        <form onSubmit={submit}>
          <label>
            Coupon code (optional)
            <input
              data-testid="coupon"
              value={coupon}
              onChange={(e) => setCoupon(e.target.value)}
              placeholder="e.g. WELCOME10"
            />
          </label>
          <label>
            Card token
            <input
              data-testid="card"
              value={card}
              onChange={(e) => setCard(e.target.value)}
              required
            />
          </label>
          <p style={{ fontSize: "0.78rem", color: "var(--gray-400)", marginTop: "-0.5rem" }}>
            Use <code>tok_decline_card</code> to test a declined payment.
          </p>
          <button
            type="submit"
            data-testid="place-order"
            disabled={loading}
            style={{ marginTop: "0.5rem", width: "100%", padding: "0.75rem" }}
          >
            {loading ? "Placing order…" : "Place order"}
          </button>
          {err && <p data-testid="checkout-error" role="alert">{err}</p>}
        </form>
      </div>
    </section>
  );
}
