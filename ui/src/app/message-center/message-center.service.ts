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

  updateSubscribe(subscribable): Observable<any> {
    return this.httpClient.post<any>(this.baseUrl + 'subscribe/' + subscribable.id + '/', subscribable).pipe(
      catchError(error => throwError(error))
    );
  }

  listUserReceiver(): Observable<any> {
    return this.httpClient.get<any>(this.baseUrl + 'receiver').pipe(
      catchError(error => throwError(error))
    );
  }

  updateUserReceiver(receiver): Observable<any> {
    return this.httpClient.post<any>(this.baseUrl + 'receiver/' + receiver.id + '/', receiver).pipe(
      catchError(error => throwError(error))
    );
  }
}
