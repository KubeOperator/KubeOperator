import {
    HttpEvent,
    HttpInterceptor,
    HttpHandler,
    HttpRequest
} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Injectable} from '@angular/core';
import { Md5 } from 'ts-md5';

@Injectable()
export class SessionInterceptor implements HttpInterceptor {

    intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
        if (req.method !== 'GET') {
            req = req.clone({ headers: req.headers.set('X-CSRF-TOKEN', this.encrypt())});
        }

        if (req.url.startsWith('/api')) {
            const currentLanguage = localStorage.getItem('currentLanguage');
            const clonedRequest = req.clone({
                params: req.params.set('l', currentLanguage)
            });
            return next.handle(clonedRequest);
        }
        return next.handle(req);
    }

    encrypt() {
        var offset = new Date().getTimezoneOffset()
        var thisTime = new Date()
        thisTime.setMinutes(thisTime.getMinutes() + offset)
        let formateDay = (day) => {
          return String(day).replace(/(^\d{1}$)/,'0$1')
        }
        var kk = formateDay(thisTime.getMonth() + 1) + "-" + formateDay(thisTime.getDate()) + " " + formateDay(thisTime.getHours()) + ":" + formateDay(thisTime.getMinutes()) + ":" + formateDay(thisTime.getSeconds())
        
        return Md5.hashStr("kubeoperator" + kk)
    }
}
