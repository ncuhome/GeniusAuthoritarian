export function numeral(num: number, digital = 1): string {
  if (num < 1000) {
    return String(num);
  }
  if (num < 1000 * 1000) {
    num = num / 1000;
    let value: string;
    if (num < 1000 * 100) {
      value = num.toFixed(digital);
    } else {
      value = num.toFixed(0);
    }
    return value + "k";
  }
  return (num / 1000 / 1000).toFixed(digital) + "m";
}
