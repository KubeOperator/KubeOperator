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

  listUserMessageByPage(limit: number, page: number, type: string, readStatus: string, level: string): Observable<any> {
    return this.httpClient.get<any>(this.baseUrl + 'userMessage?limit=' + limit + '&page=' +
      page + '&type=' + type + '&readStatus=' + readStatus + '&level=' + level).pipe(
      catchError(error => throwError(error))
    );
  }

  updateUserMessage(userMessage): Observable<any> {
    return this.httpClient.post<any>(this.baseUrl + 'userMessage/' + userMessage.id + '/', userMessage).pipe(
      catchError(error => throwError(error))
    );
  }

  updateAllUserMessage(): Observable<any> {
    return this.httpClient.post<any>(this.baseUrl + 'userMessage/ALL/', {}).pipe(
      catchError(error => throwError(error))
    );
  }

  unReadMessage(): Observable<any> {
    return this.httpClient.get<any>(this.baseUrl + 'userMessage/unread/').pipe(
      catchError(error => throwError(error))
    );
  }
}
