import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Ceph} from './ceph';


@Injectable({
  providedIn: 'root'
})
export class CephService {

  base_url = '/api/v1/storage/ceph/';

  constructor(private http: HttpClient) {
  }

  list(): Observable<Ceph[]> {
    return this.http.get<Ceph[]>(this.base_url);
  }

  delete(name: string): Observable<any> {
    return this.http.delete<any>(this.base_url + name + '/');
  }

  create(item: Ceph): Observable<Ceph> {
    return this.http.post<Ceph>(this.base_url, item);
  }
}
