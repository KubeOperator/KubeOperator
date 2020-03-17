import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http' ;
import {Observable, throwError} from 'rxjs';
import {catchError} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})

export class MessageCenterService {

  constructor(private httpClient: HttpClient) {
  }

  baseUrl = '/api/v1/notification/';


  listSubscribe(): Observable<any> {
    return this.httpClient.get<any>(this.baseUrl + 'subscribe/').pipe(
      catchError(error => throwError(error))
    );
  }
}
