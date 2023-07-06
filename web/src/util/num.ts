export function numeral(num: number, digital = 1): string {
  if (num < 1000) {
    return String(num);
  }
  if (num < 1000 * 1000) {
    return (num / 1000).toFixed(digital) + "k";
  }
  return (num / 1000 / 1000).toFixed(digital) + "m";
}
