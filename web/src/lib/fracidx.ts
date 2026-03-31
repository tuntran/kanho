import { generateKeyBetween, generateNKeysBetween } from "fractional-indexing";

export { generateKeyBetween, generateNKeysBetween };

/** Generate a key at the start (before all existing keys). */
export function keyStart(): string {
  return generateKeyBetween(null, null);
}

/** Generate a key after the last key. */
export function keyAfter(last: string): string {
  return generateKeyBetween(last, null);
}

/** Generate a key between two keys. */
export function keyBetween(a: string | null, b: string | null): string {
  return generateKeyBetween(a, b);
}
