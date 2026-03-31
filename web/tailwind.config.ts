import type { Config } from "tailwindcss";

export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        bg: "#0F1117",
        surface: "#1C1F2A",
        "surface-2": "#252836",
        "surface-raised": "#242838",
        border: "#2D3144",
        "border-focus": "#6366F1",
        primary: "#6366F1",
        "primary-hover": "#4F46E5",
        accent: "#A78BFA",
        success: "#22C55E",
        warning: "#F59E0B",
        danger: "#EF4444",
        info: "#38BDF8",
        "text-1": "#F1F5F9",
        "text-2": "#94A3B8",
        "text-3": "#475569",
      },
      fontFamily: {
        sans: ["Inter", "system-ui", "sans-serif"],
        mono: ["JetBrains Mono", "monospace"],
      },
      borderRadius: {
        card: "8px",
        panel: "12px",
      },
      boxShadow: {
        card: "0 1px 4px rgba(0,0,0,0.3)",
        "card-hover": "0 4px 16px rgba(0,0,0,0.4)",
        "card-drag": "0 8px 32px rgba(99,102,241,0.3)",
        drawer: "-8px 0 32px rgba(0,0,0,0.4)",
      },
      keyframes: {
        "pulse-dot": {
          "0%, 100%": {
            opacity: "1",
            boxShadow: "0 0 0 0 rgba(34,197,94,0.4)",
          },
          "50%": {
            opacity: "0.8",
            boxShadow: "0 0 0 6px rgba(34,197,94,0)",
          },
        },
      },
      animation: {
        "pulse-dot": "pulse-dot 2s infinite",
      },
    },
  },
  plugins: [],
} satisfies Config;
