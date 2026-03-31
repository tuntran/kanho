import { useState } from "react";
import type { FormEvent } from "react";
import { Link } from "react-router-dom";
import { useAuth } from "@/hooks/use-auth";

export function RegisterPage() {
  const { register } = useAuth();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");

    if (password.length < 8) {
      setError("Password must be at least 8 characters");
      return;
    }

    setLoading(true);
    try {
      await register(email, password, name);
    } catch (err: unknown) {
      if (typeof err === "object" && err !== null && "response" in err) {
        const resp = err as { response?: { data?: { error?: string } } };
        setError(resp.response?.data?.error ?? "Registration failed");
      } else {
        setError(err instanceof Error ? err.message : "Registration failed");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-bg px-4">
      <div className="w-full max-w-[400px]">
        {/* Logo */}
        <div className="mb-8 text-center">
          <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-panel bg-primary/10">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" className="text-primary">
              <rect x="3" y="3" width="7" height="7" rx="1.5" fill="currentColor" opacity="0.9" />
              <rect x="14" y="3" width="7" height="7" rx="1.5" fill="currentColor" opacity="0.6" />
              <rect x="3" y="14" width="7" height="7" rx="1.5" fill="currentColor" opacity="0.6" />
              <rect x="14" y="14" width="7" height="7" rx="1.5" fill="currentColor" opacity="0.3" />
            </svg>
          </div>
          <h1 className="text-xl font-bold text-text-1">Create your account</h1>
          <p className="mt-1 text-sm text-text-2">Get started with Kanho</p>
        </div>

        {/* Card */}
        <div className="rounded-panel border border-border bg-surface p-6 shadow-card">
          {error && (
            <div className="mb-4 rounded-card border border-danger/20 bg-danger/10 px-3 py-2.5 text-sm text-danger">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label htmlFor="name" className="mb-1.5 block text-sm font-medium text-text-2">
                Full name
              </label>
              <input
                id="name"
                type="text"
                required
                autoComplete="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="h-10 w-full rounded-card border border-border bg-bg px-3 text-sm text-text-1 placeholder-text-3 outline-none transition-colors focus:border-border-focus focus:ring-1 focus:ring-primary/30"
                placeholder="Your full name"
              />
            </div>

            <div>
              <label htmlFor="email" className="mb-1.5 block text-sm font-medium text-text-2">
                Email
              </label>
              <input
                id="email"
                type="email"
                required
                autoComplete="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="h-10 w-full rounded-card border border-border bg-bg px-3 text-sm text-text-1 placeholder-text-3 outline-none transition-colors focus:border-border-focus focus:ring-1 focus:ring-primary/30"
                placeholder="you@example.com"
              />
            </div>

            <div>
              <label htmlFor="password" className="mb-1.5 block text-sm font-medium text-text-2">
                Password
              </label>
              <input
                id="password"
                type="password"
                required
                minLength={8}
                autoComplete="new-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="h-10 w-full rounded-card border border-border bg-bg px-3 text-sm text-text-1 placeholder-text-3 outline-none transition-colors focus:border-border-focus focus:ring-1 focus:ring-primary/30"
                placeholder="Min 8 characters"
              />
              <p className="mt-1.5 text-xs text-text-3">Must be at least 8 characters long</p>
            </div>

            <button
              type="submit"
              disabled={loading}
              className="h-10 w-full rounded-card bg-primary text-sm font-semibold text-white transition-all hover:-translate-y-px hover:bg-primary-hover hover:brightness-110 active:translate-y-0 active:brightness-95 disabled:pointer-events-none disabled:opacity-50"
            >
              {loading ? (
                <span className="inline-flex items-center gap-2">
                  <svg className="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
                    <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="3" strokeLinecap="round" className="opacity-25" />
                    <path d="M4 12a8 8 0 018-8" stroke="currentColor" strokeWidth="3" strokeLinecap="round" />
                  </svg>
                  Creating account...
                </span>
              ) : (
                "Create account"
              )}
            </button>
          </form>
        </div>

        {/* Footer link */}
        <p className="mt-6 text-center text-sm text-text-2">
          Already have an account?{" "}
          <Link to="/login" className="font-medium text-accent transition-colors hover:text-primary">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  );
}
