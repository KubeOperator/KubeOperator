import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {CloudZone} from './cloud';

@Injectable({
  providedIn: 'root'
})
export class CloudService {
  regionUrl = '/api/v1/cloud/region/';
  zoneUrl = '/api/v1/cloud/{region}/zone/';

  constructor(private http: HttpClient) {
  }

  listRegion(vars: any): Observable<string[]> {
    return this.http.post<string[]>(this.regionUrl, vars);
  }

  listZone(region: string): Observable<CloudZone[]> {
    return this.http.get<CloudZone[]>(this.zoneUrl.replace('{region}', region));
  }
}
