import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { orders } from "../lib/api";

export default function Checkout() {
  const [coupon, setCoupon] = useState("");
  const [card, setCard] = useState("tok_test_visa");
  const [err, setErr] = useState("");
  const nav = useNavigate();

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setErr("");
    try {
      const o = await orders.checkout(coupon, card, crypto.randomUUID());
      nav(`/orders/${o.id}`);
    } catch (e: any) {
      setErr(e.message);
    }
  }

  return (
    <section data-testid="checkout-page">
      <h1>Checkout</h1>
      <form onSubmit={submit}>
        <label>
          Coupon
          <input data-testid="coupon" value={coupon} onChange={(e) => setCoupon(e.target.value)} />
        </label>
        <label>
          Card token
          <input data-testid="card" value={card} onChange={(e) => setCard(e.target.value)} required />
        </label>
        <button type="submit" data-testid="place-order">Place order</button>
        {err && <p data-testid="checkout-error" role="alert">{err}</p>}
      </form>
    </section>
  );
}
