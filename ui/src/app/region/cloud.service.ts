import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class CloudService {
  regionUrl = '/api/v1/cloud/region';

  constructor(private http: HttpClient) {
  }

  listRegion(vars: any): Observable<string[]> {
    return this.http.post<string[]>(this.regionUrl, vars);
  }
}
