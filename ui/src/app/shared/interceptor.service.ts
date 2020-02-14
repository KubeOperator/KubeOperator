import {Injectable} from '@angular/core';
import {HttpErrorResponse, HttpEvent, HttpHandler, HttpInterceptor, HttpRequest} from '@angular/common/http';
import {Observable, throwError} from 'rxjs';
import {catchError, map} from 'rxjs/operators';
import {JwtHelperService} from '@auth0/angular-jwt';
import {SessionService} from './session.service';
import {CommonAlertService} from '../base/header/common-alert.service';
import {AlertLevels} from '../base/header/components/common-alert/alert';


@Injectable()
export class InterceptorService implements HttpInterceptor {
  constructor(private session: SessionService, private alert: CommonAlertService) {
  }

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    let token = null;
    const sessionUser = JSON.parse(localStorage.getItem('current_user'));
    if (sessionUser) {
      token = sessionUser.token;
      const helper = new JwtHelperService();
      const expirationDate = helper.getTokenExpirationDate(token);
      const now = new Date();
      if (now.getTime() < expirationDate.getTime() && expirationDate.getTime() - now.getTime() <= 1000 * 10 * 60) {
        this.session.refreshToken(token).subscribe((data) => {
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
        let msg = '无可用消息！';
        if (error.status === 403) {
          msg = '权限不允许此操作或者Session过期！';
          this.alert.showAlert(msg, AlertLevels.ERROR);
        }
        return throwError(error);
      }));
  }
}

