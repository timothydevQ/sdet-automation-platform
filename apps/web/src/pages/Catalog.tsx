import React, { useEffect, useState } from "react";
import { catalog, cart } from "../lib/api";

const EMOJI: Record<string, string> = {
  electronics: "💻",
  accessories: "🔌",
  merch: "🛍️",
  general: "📦",
};

export default function Catalog() {
  const [items, setItems] = useState<any[]>([]);
  const [msg, setMsg] = useState("");
  const [adding, setAdding] = useState<string | null>(null);

  useEffect(() => {
    catalog.list().then(setItems).catch(() => setItems([]));
  }, []);

  async function add(sku: string) {
    setAdding(sku);
    try {
      await cart.add(sku, 1);
      setMsg("Added to cart");
      setTimeout(() => setMsg(""), 2500);
    } catch (e: any) {
      setMsg(e.message);
    } finally {
      setAdding(null);
    }
  }

  return (
    <section data-testid="catalog-page">
      <h1>Shop</h1>
      <p className="catalog-subtitle">{items.length} products</p>
      {msg && <p data-testid="catalog-msg">{msg}</p>}
      <div className="product-grid">
        {items.map((p) => (
          <div key={p.sku} className="product-card" data-testid={`product-${p.sku}`}>
            <div className="product-img">{EMOJI[p.category] ?? "📦"}</div>
            <div className="product-info">
              <div className="product-category">{p.category}</div>
              <div className="product-name">{p.name}</div>
              <div className="product-footer">
                <div>
                  <div className="product-price">${(p.price_cents / 100).toFixed(2)}</div>
                  <div className="product-stock">
                    {p.stock > 10 ? "In stock" : p.stock > 0 ? `Only ${p.stock} left` : "Out of stock"}
                  </div>
                </div>
                <button
                  onClick={() => add(p.sku)}
                  data-testid={`add-${p.sku}`}
                  disabled={adding === p.sku || p.stock === 0}
                >
                  {adding === p.sku ? "Adding…" : "Add"}
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}
