// Will find the smallest integer k in [0, n) where f(k)=true
// This assumes if f(i)=true then f(i+1)=true
// If none are true returns n
export function BinarySearch(f: (index: number) => boolean, n: number): number {
  if (n <= 0 || f(0)) {
    return 0;
  }

  let l = 0; // we know f(l)=false
  let r = n; // not inclusive of r. We know f(r)=true
  while (r - l > 1) {
    // if l=0 and r=1 Math.ceil would fail
    const mid = Math.floor((l + r) / 2);
    if (f(mid)) {
      r = mid;
    } else {
      l = mid;
    }
  }

  return r;
}
