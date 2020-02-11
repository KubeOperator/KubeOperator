import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, throwError} from 'rxjs';
import {catchError} from 'rxjs/operators';
import {Item} from './item';

@Injectable({
  providedIn: 'root'
})
export class ItemService {
  baseUrl = '/api/v1/items/';

  constructor(private http: HttpClient) {
  }

  listItem(): Observable<Item[]> {
    return this.http.get<Item[]>(this.baseUrl);
  }

  getItem(itemName: string): Observable<Item> {
    return this.http.get<Item>(this.baseUrl + itemName + '/').pipe(
      catchError(err => throwError(err))
    );
  }

  deleteItem(itemName): Observable<any> {
    return this.http.delete(this.baseUrl + itemName + '/').pipe(
      catchError(error => throwError(error))
    );
  }

  createItem(item: Item): Observable<any> {
    return this.http.post(this.baseUrl, item).pipe(
      catchError(error => throwError(error))
    );
  }
}
