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
        let formateDay = (day) => {
            return String(day).replace(/(^\d{1}$)/,'0$1')
        }
        let dateNow = new Date()
        let month = formateDay(dateNow.getUTCMonth() + 1)
        let date = formateDay(dateNow.getUTCDate())
        let hour = formateDay(dateNow.getUTCHours())
        let minute = formateDay(dateNow.getUTCMinutes())
        let second = formateDay(dateNow.getUTCSeconds())
        var kk = month + "-" + date + " " + hour + ":" + minute + ":" + second
        console.log(kk, Md5.hashStr("kubeoperator" + kk))
        
        return Md5.hashStr("kubeoperator" + kk)
    }
}
