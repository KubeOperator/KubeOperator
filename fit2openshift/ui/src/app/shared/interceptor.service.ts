import {Injectable} from '@angular/core';
import {HttpEvent, HttpHandler, HttpInterceptor, HttpRequest, HttpResponse} from '@angular/common/http';
import {Observable} from 'rxjs';
import {mergeMap} from 'rxjs/operators';
import {MessageService} from '../base/message.service';
import {MessageLevels} from '../base/message/message-level';

@Injectable()
export class InterceptorService implements HttpInterceptor {

  constructor(private messageService: MessageService) {
  }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    let token = null;
    const sessionUser = JSON.parse(localStorage.getItem('current_user'));
    if (sessionUser) {
      token = sessionUser.token;
    }
    if (token) {
      req = req.clone({
        setHeaders: {
          Authorization: `JWT ${token}`
        }
      });
    }

    return next.handle(req).pipe(
      mergeMap((event: any) => {
        if (event instanceof HttpResponse && event.status !== 200) {
          this.messageService.announceMessage(event.statusText, MessageLevels.ERROR);
        }
        return Observable.create(observer => observer.next(event));
      })
    );
  }


}
