import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { auth } from "../lib/api";

export default function Login() {
  const [mode, setMode] = useState<"login" | "register">("login");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [err, setErr] = useState("");
  const [loading, setLoading] = useState(false);
  const nav = useNavigate();

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setErr("");
    setLoading(true);
    try {
      const r = mode === "login"
        ? await auth.login(email, password)
        : await auth.register(email, password);
      localStorage.setItem("token", r.token);
      localStorage.setItem("role", r.role);
      nav("/");
    } catch {
      setErr(mode === "login" ? "Invalid email or password" : "Registration failed — account may already exist");
    } finally {
      setLoading(false);
    }
  }

  return (
    <section data-testid="login-page">
      <div className="form-card">
        <div style={{ display: "flex", borderBottom: "1px solid var(--gray-200)", marginBottom: "1.75rem" }}>
          {(["login", "register"] as const).map((m) => (
            <button
              key={m}
              onClick={() => { setMode(m); setErr(""); }}
              style={{
                flex: 1,
                background: "none",
                border: "none",
                borderBottom: mode === m ? "2px solid var(--black)" : "2px solid transparent",
                borderRadius: 0,
                padding: "0.75rem",
                fontWeight: mode === m ? 700 : 400,
                color: mode === m ? "var(--black)" : "var(--gray-400)",
                cursor: "pointer",
                fontSize: "0.9rem",
                textTransform: "capitalize",
                letterSpacing: "0.02em",
              }}
            >
              {m === "login" ? "Sign in" : "Create account"}
            </button>
          ))}
        </div>

        <form onSubmit={submit}>
          <label>
            Email
            <input
              data-testid="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
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
              placeholder="••••••••"
              minLength={6}
              required
            />
          </label>
          {mode === "register" && (
            <p style={{ fontSize: "0.78rem", color: "var(--gray-400)", marginTop: "-0.5rem" }}>
              Tip: use an email ending in <code>@admin.local</code> for admin access.
            </p>
          )}
          <button
            data-testid="submit"
            type="submit"
            disabled={loading}
            style={{ marginTop: "0.5rem", width: "100%", padding: "0.75rem" }}
          >
            {loading ? "Please wait…" : mode === "login" ? "Sign in" : "Create account"}
          </button>
          {err && <p data-testid="error" role="alert">{err}</p>}
        </form>
      </div>
    </section>
  );
}
