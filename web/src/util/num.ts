export function numeral(num: number): string {
  if (num < 1000) {
    return String(num);
  }
  if (num < 1000 * 1000) {
    return (num / 1000).toFixed(1) + "k";
  }
  return (num / 1000 / 1000).toFixed(1) + "m";
}
