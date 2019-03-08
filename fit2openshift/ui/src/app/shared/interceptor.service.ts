import {Injectable} from '@angular/core';
import {HttpErrorResponse, HttpEvent, HttpHandler, HttpInterceptor, HttpRequest, HttpResponse} from '@angular/common/http';
import {Observable, of, throwError} from 'rxjs';
import {catchError, map} from 'rxjs/operators';
import {MessageService} from '../base/message.service';


@Injectable()
export class InterceptorService implements HttpInterceptor {
  constructor(private message: MessageService) {
  }

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    let token = null;
    const sessionUser = JSON.parse(localStorage.getItem('current_user'));
    if (sessionUser) {
      token = sessionUser.token;
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
        this.message.messagesQueue.next(error.error.reason);
        let data = {}
        ;
        data = {
          reason: error && error.error.reason ? error.error.reason : '',
          status: error.status
        };
        return throwError(data);
      }));
  }
}


// export class InterceptorService implements HttpInterceptor {
//
//   constructor(private messageService: MessageService) {
//   }
//
//   intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
//     let token = null;
//     const sessionUser = JSON.parse(localStorage.getItem('current_user'));
//     if (sessionUser) {
//       token = sessionUser.token;
//     }
//     if (token) {
//       req = req.clone({
//         setHeaders: {
//           Authorization: `JWT ${token}`
//         }
//       });
//     }
//     //
//     // return next.handle(req).pipe(tap(data => console.log(data)));
//
//     return next.handle(req).pipe(mergeMap((event: any) => {
//       if (event instanceof HttpResponse && event.status === 200) {
//         return this.handleData(event);
//       }
//       return of(event);
//     }));
//   }
//
//
//   handleData(event: HttpResponse<any> | HttpErrorResponse): Observable<any> {
//     switch (event.status) {
//       case 200:
//         console.log('ok');
//         break;
//       case 401:
//         break;
//       case 404:
//         break;
//       case 504:
//         break;
//       default:
//         return of(event);
//     }
//   }
// }
