import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { cart } from "../lib/api";

export default function Cart() {
  const [items, setItems] = useState<any[]>([]);

  useEffect(() => {
    cart.get().then(setItems).catch(() => setItems([]));
  }, []);

  const total = items.reduce((s, i) => s + i.price_cents * i.qty, 0);

  return (
    <section data-testid="cart-page">
      <h1>Cart</h1>
      <ul>
        {items.map((i) => (
          <li key={i.sku} data-testid={`cart-item-${i.sku}`}>
            {i.name} × {i.qty}
            <button onClick={() => cart.remove(i.sku).then(() => cart.get().then(setItems))}>
              Remove
            </button>
          </li>
        ))}
      </ul>
      <p data-testid="cart-total">Total: ${(total / 100).toFixed(2)}</p>
      <Link to="/checkout">
        <button data-testid="checkout-btn" disabled={items.length === 0}>Checkout</button>
      </Link>
    </section>
  );
}
