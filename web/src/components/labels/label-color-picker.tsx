import { useState } from "react";

const PRESET_COLORS = [
  "#ef4444", "#f97316", "#eab308", "#22c55e",
  "#14b8a6", "#3b82f6", "#6366f1", "#a855f7",
  "#ec4899", "#f43f5e", "#78716c", "#64748b",
  "#6b7280", "#059669", "#0891b2", "#7c3aed",
];

interface LabelColorPickerProps {
  value: string;
  onChange: (color: string) => void;
}

export function LabelColorPicker({ value, onChange }: LabelColorPickerProps) {
  const [hexInput, setHexInput] = useState(value);

  const handleHexChange = (raw: string) => {
    setHexInput(raw);
    const hex = raw.startsWith("#") ? raw : `#${raw}`;
    if (/^#[0-9a-fA-F]{6}$/.test(hex)) {
      onChange(hex);
    }
  };

  return (
    <div className="space-y-2">
      <div className="grid grid-cols-8 gap-1.5">
        {PRESET_COLORS.map((color) => (
          <button
            key={color}
            type="button"
            onClick={() => {
              onChange(color);
              setHexInput(color);
            }}
            className={`h-6 w-6 rounded-full border-2 transition-transform hover:scale-110 ${
              value === color ? "border-text-1 scale-110" : "border-transparent"
            }`}
            style={{ backgroundColor: color }}
          />
        ))}
      </div>
      <input
        type="text"
        value={hexInput}
        onChange={(e) => handleHexChange(e.target.value)}
        placeholder="#hex"
        className="w-full rounded-md border border-border bg-surface-2 px-2 py-1 text-xs text-text-1 placeholder:text-text-2 focus:border-accent focus:outline-none"
      />
    </div>
  );
}
