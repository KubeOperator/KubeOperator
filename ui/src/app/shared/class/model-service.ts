import {Observable} from 'rxjs';
import {Page} from './page';
import {HttpClient} from '@angular/common/http';

export class ModelService<T> {
  protected baseUrl = '';

  constructor(protected http: HttpClient) {
  }

  list(page: number, size: number): Observable<Page<T>> {
    return this.http.get<Page<T>>(`${this.baseUrl}?page=${page}&size=${size}`);
  }

  listAll(): Observable<T[]> {
    return this.http.get<T[]>(this.baseUrl);
  }

  get(pk: string): Observable<T> {
    return this.http.get<T>(`${this.baseUrl}${pk}/`);
  }

  create(item: T): Observable<T> {
    return this.http.post<T>(`${this.baseUrl}`, item);
  }

  delete(pk: string): Observable<any> {
    return this.http.delete<any>(`${this.baseUrl}${pk}/`);
  }

  update(item: T, pk: string): Observable<T> {
    return this.http.patch<T>(`${this.baseUrl}${pk}/`, item);
  }

}
