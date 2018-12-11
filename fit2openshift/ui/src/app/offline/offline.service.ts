import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, throwError} from 'rxjs';
import {Offline} from './Offline';
import {catchError} from 'rxjs/operators';

const offlineUrl = '/api/v1/offline/';

@Injectable()
export class OfflineService {

  constructor(private http: HttpClient) {
  }

  listOfflines(): Observable<Offline[]> {
    return this.http.get<Offline[]>(offlineUrl).pipe(
      catchError(error => throwError(error))
    );
  }

  createOffline(offline: Offline): Observable<Offline> {
    return this.http.post<Offline>(offlineUrl, offline).pipe(
      catchError(error => throwError(error))
    );
  }

  getOffline(offlineId: string): Observable<Offline> {
    return this.http.get<Offline>(`${offlineUrl}/${offlineId}`).pipe(
      catchError(error => throwError(error))
    );
  }

  deleteOffline(offlineId: string): Observable<any> {
    return this.http.delete(offlineId).pipe(
      catchError(error => throwError(error))
    );
  }

}
