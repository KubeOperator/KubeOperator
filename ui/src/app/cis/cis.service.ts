import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, throwError} from 'rxjs';
import {catchError} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class CisService {


  constructor(private httpClient: HttpClient) {
  }


  baseUrl = '/api/v1/cisLog/';
  cisUrl = '/api/v1/cluster/{cluster_id}/cisLog/';


  listCis(cluster_id: string): Observable<any> {
    return this.httpClient.get<any>(this.cisUrl.replace('{cluster_id}', cluster_id)).pipe(
      catchError(error => throwError(error))
    );
  }

  deleteCis(name: string): Observable<any> {
    return this.httpClient.delete<any>(this.baseUrl + name + '/').pipe(
      catchError(error => throwError(error))
    );
  }

  runCis(cluster_id: string): Observable<any> {
    return this.httpClient.post<any>((this.cisUrl + 'run/').replace('{cluster_id}', cluster_id), {}).pipe(
      catchError(error => throwError(error))
    );
  }
}
