import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {NfsStorage} from './nfs';
import {Observable} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class NfsService {

  baseUrl = '/api/v1/storage/nfs/';

  constructor(private http: HttpClient) {
  }

  list(): Observable<NfsStorage[]> {
    return this.http.get<NfsStorage[]>(this.baseUrl);
  }

  create(item: NfsStorage): Observable<NfsStorage> {
    return this.http.post<NfsStorage>(this.baseUrl, item);
  }

  delete(name: string): Observable<NfsStorage> {
    return this.http.delete<NfsStorage>(this.baseUrl + name + '/');
  }

}
