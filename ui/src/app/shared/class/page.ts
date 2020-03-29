export class Page<T> {
  count: number;
  next: number;
  previous: number;
  results: T[] = [];
}
