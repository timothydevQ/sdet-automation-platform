import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { cart } from "../lib/api";

export default function Cart() {
  const [items, setItems] = useState<any[]>([]);
  const nav = useNavigate();

  async function load() {
    cart.get().then(setItems).catch(() => setItems([]));
  }

  useEffect(() => { load(); }, []);

  const total = items.reduce((s, i) => s + i.price_cents * i.qty, 0);

  if (items.length === 0) {
    return (
      <section data-testid="cart-page">
        <h1>Cart</h1>
        <div className="cart-empty">
          <p style={{ fontSize: "2.5rem", marginBottom: "1rem" }}>🛒</p>
          <p>Your cart is empty.</p>
          <button onClick={() => nav("/")} style={{ marginTop: "1.5rem" }}>Browse products</button>
        </div>
      </section>
    );
  }

  return (
    <section data-testid="cart-page">
      <h1>Your cart</h1>
      <table className="cart-table">
        <thead>
          <tr>
            <th>Product</th>
            <th>Price</th>
            <th>Qty</th>
            <th>Subtotal</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {items.map((i) => (
            <tr key={i.sku} data-testid={`cart-item-${i.sku}`}>
              <td style={{ fontWeight: 500 }}>{i.name}</td>
              <td>${(i.price_cents / 100).toFixed(2)}</td>
              <td>{i.qty}</td>
              <td>${(i.price_cents * i.qty / 100).toFixed(2)}</td>
              <td>
                <button
                  className="btn-ghost"
                  onClick={() => cart.remove(i.sku).then(load)}
                >
                  Remove
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <div className="cart-summary">
        <span className="cart-total-label">Total</span>
        <span className="cart-total-value" data-testid="cart-total">${(total / 100).toFixed(2)}</span>
      </div>
      <div className="cart-actions">
        <button data-testid="checkout-btn" onClick={() => nav("/checkout")}>
          Proceed to checkout →
        </button>
      </div>
    </section>
  );
}
