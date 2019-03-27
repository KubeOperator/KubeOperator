import {Injectable} from '@angular/core';
import {HttpErrorResponse, HttpEvent, HttpHandler, HttpInterceptor, HttpRequest} from '@angular/common/http';
import {Observable, throwError} from 'rxjs';
import {catchError, map} from 'rxjs/operators';
import {MessageService} from '../base/message.service';
import {MessageLevels} from '../base/message/message-level';
import {JwtHelperService} from '@auth0/angular-jwt';
import {SessionService} from './session.service';


@Injectable()
export class InterceptorService implements HttpInterceptor {
  constructor(private message: MessageService, private session: SessionService) {
  }

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    let token = null;
    const sessionUser = JSON.parse(localStorage.getItem('current_user'));
    if (sessionUser) {
      token = sessionUser.token;
      const helper = new JwtHelperService();
      const expirationDate = helper.getTokenExpirationDate(token);
      const now = new Date();
      if (expirationDate.getTime() - now.getTime() <= 1000 * 10 * 60) {
        this.session.refreshToken().subscribe((data) => {
          this.session.cacheToken(data);
        });
      }
    }
    if (token) {
      request = request.clone({
        setHeaders: {
          Authorization: `JWT ${token}`
        }
      });
    }
    return next.handle(request).pipe(
      map((event: HttpEvent<any>) => {
        return event;
      }),
      catchError((error: HttpErrorResponse) => {
        console.log(error);
        let msg = '无可用消息！';
        if (error.status === 403) {
          msg = '权限不允许此操作或者Session过期！';
          this.message.announceMessage(msg, MessageLevels.ERROR);
        }
        return throwError(error);
      }));
  }
}

