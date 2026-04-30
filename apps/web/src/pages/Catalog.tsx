import React, { useEffect, useState } from "react";
import { catalog, cart } from "../lib/api";

export default function Catalog() {
  const [items, setItems] = useState<any[]>([]);
  const [msg, setMsg] = useState("");

  useEffect(() => {
    catalog.list().then(setItems).catch(() => setItems([]));
  }, []);

  async function add(sku: string) {
    try {
      await cart.add(sku, 1);
      setMsg(`Added ${sku}`);
    } catch (e: any) {
      setMsg(e.message);
    }
  }

  return (
    <section data-testid="catalog-page">
      <h1>Catalog</h1>
      {msg && <p data-testid="catalog-msg">{msg}</p>}
      <ul>
        {items.map((p) => (
          <li key={p.sku} data-testid={`product-${p.sku}`}>
            <strong>{p.name}</strong> — ${(p.price_cents / 100).toFixed(2)}
            <button onClick={() => add(p.sku)} data-testid={`add-${p.sku}`}>
              Add to cart
            </button>
          </li>
        ))}
      </ul>
    </section>
  );
}
