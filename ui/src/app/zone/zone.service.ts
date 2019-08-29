import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Zone} from './zone';

@Injectable({
  providedIn: 'root'
})
export class ZoneService {
  baseUrl = '/api/v1/zones/';

  constructor(private http: HttpClient) {
  }

  listZones(): Observable<Zone[]> {
    return this.http.get<Zone[]>(this.baseUrl);
  }

  createZones(item: Zone): Observable<Zone> {
    return this.http.post<Zone>(this.baseUrl, item);
  }

  getZone(name: string): Observable<Zone> {
    return this.http.get<Zone>(this.baseUrl + name + '/');
  }

  deleteZone(name: string): Observable<Zone> {
    return this.http.delete<Zone>(this.baseUrl + name + '/');
  }
}
