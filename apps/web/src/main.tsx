import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route, Link, useNavigate } from "react-router-dom";
import Login from "./pages/Login";
import Catalog from "./pages/Catalog";
import Cart from "./pages/Cart";
import Checkout from "./pages/Checkout";
import AdminOrders from "./pages/AdminOrders";
import "./styles.css";

function Nav() {
  const nav = useNavigate();
  function logout() {
    localStorage.removeItem("token");
    localStorage.removeItem("role");
    nav("/login");
  }
  const token = localStorage.getItem("token");
  const role = localStorage.getItem("role");

  return (
    <nav data-testid="nav">
      <Link to="/">SDET Shop</Link>
      <Link to="/">Catalog</Link>
      <Link to="/cart">Cart</Link>
      {role === "admin" && <Link to="/admin/orders">Admin</Link>}
      {token
        ? <button onClick={logout} style={{ background: "transparent", color: "var(--gray-400)", fontSize: "0.82rem", padding: "0.3rem 0.7rem", border: "1px solid #444" }}>Sign out</button>
        : <Link to="/login">Sign in</Link>
      }
    </nav>
  );
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <BrowserRouter>
    <Nav />
    <main>
      <Routes>
        <Route path="/" element={<Catalog />} />
        <Route path="/login" element={<Login />} />
        <Route path="/cart" element={<Cart />} />
        <Route path="/checkout" element={<Checkout />} />
        <Route path="/admin/orders" element={<AdminOrders />} />
      </Routes>
    </main>
  </BrowserRouter>
);
