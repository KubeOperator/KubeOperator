import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Region} from './region';

@Injectable({
  providedIn: 'root'
})
export class RegionService {


  constructor(private http: HttpClient) {
  }

  baseUrl = '/api/v1/regions/';

  listRegion(): Observable<Region[]> {
    return this.http.get<Region[]>(this.baseUrl);
  }

  getRegion(name: string): Observable<Region> {
    return this.http.get<Region>(this.baseUrl + name + '/');
  }


  createRegion(item: Region): Observable<Region> {
    return this.http.post<Region>(this.baseUrl, item);
  }

  updateRegion(name: string, item: Region): Observable<Region> {
    return this.http.patch<Region>(this.baseUrl + name + '/', item);
  }

  deleteRegion(name: string): Observable<Region> {
    return this.http.delete<Region>(this.baseUrl + name + '/');
  }
}
