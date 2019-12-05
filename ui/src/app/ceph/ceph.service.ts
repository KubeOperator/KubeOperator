import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';


@Injectable({
  providedIn: 'root'
})
export class CephService {

  base_url = '/api/v1/storage/ceph';

  constructor(private http: HttpClient) {
  }

  list(): Observable<any> {
    return this.http.get<any>(this.base_url);
  }

  delete(name: string): Observable<any> {
    return this.http.delete<any>(this.base_url + name + '/');
  }
}
