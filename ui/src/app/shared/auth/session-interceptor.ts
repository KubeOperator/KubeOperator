import {
    HttpEvent,
    HttpInterceptor,
    HttpHandler,
    HttpRequest
} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Injectable} from '@angular/core';

@Injectable()
export class SessionInterceptor implements HttpInterceptor {

    intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
        if (req.method !== 'GET') {
            req = req.clone({ headers: req.headers.set('X-CSRF-TOKEN', localStorage.getItem("cs"))});
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
}
