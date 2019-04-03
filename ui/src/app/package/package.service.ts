import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, throwError} from 'rxjs';
import {Package} from './package';
import {catchError} from 'rxjs/operators';

const packageUrl = '/api/v1/packages/';

@Injectable()
export class PackageService {

  constructor(private http: HttpClient) {
  }

  listPackage(): Observable<Package[]> {
    return this.http.get<Package[]>(packageUrl).pipe(
      catchError(error => throwError(error))
    );
  }

  createPackage(pak: Package): Observable<Package> {
    return this.http.post<Package>(packageUrl, pak).pipe(
      catchError(error => throwError(error))
    );
  }

  getPackage(packageName: string): Observable<Package> {
    return this.http.get<Package>(`${packageUrl}${packageName}`).pipe(
      catchError(error => throwError(error))
    );
  }

  deletePackage(packageName: string): Observable<any> {
    return this.http.delete(packageName).pipe(
      catchError(error => throwError(error))
    );
  }

}
