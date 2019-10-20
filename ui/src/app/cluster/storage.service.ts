import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Storage} from './cluster';

@Injectable({
  providedIn: 'root'
})
export class StorageService {

  baseUrl = '/api/v1/storage/{type}/';

  constructor(private http: HttpClient) {
  }

  list(type: string): Observable<Storage[]> {
    return this.http.get<Storage[]>(this.baseUrl.replace('{type}', type));
  }
}
