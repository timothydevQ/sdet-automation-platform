import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { auth } from "../lib/api";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [err, setErr] = useState("");
  const nav = useNavigate();

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setErr("");
    try {
      const r = await auth.login(email, password);
      localStorage.setItem("token", r.token);
      localStorage.setItem("role", r.role);
      nav("/");
    } catch (e: any) {
      setErr(e.message);
    }
  }

  return (
    <section data-testid="login-page">
      <h1>Sign in</h1>
      <form onSubmit={submit}>
        <label>
          Email
          <input
            data-testid="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </label>
        <label>
          Password
          <input
            data-testid="password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </label>
        <button data-testid="submit" type="submit">Login</button>
        {err && <p data-testid="error" role="alert">{err}</p>}
      </form>
    </section>
  );
}
