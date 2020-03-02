import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, CanActivate, CanActivateChild, Router, RouterStateSnapshot} from '@angular/router';
import {Observable} from 'rxjs';
import {SessionService} from '../session.service';
import {CommonRoutes} from '../shared.const';

@Injectable()
export class AuthUserActiveService implements CanActivate, CanActivateChild {

  constructor(private authService: SessionService, private router: Router) {
  }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> {
    return new Observable<boolean>((observer) => {
      this.isLogin().subscribe(() => {
        this.authService.getProfile().subscribe(d => {
          const profile = JSON.parse(localStorage.getItem('profile'));
          d.user.token = profile.user.token;
          this.authService.cacheProfile(d);
          observer.next(true);
          observer.complete();
        });
      }, () => {
        observer.next(false);
        observer.complete();
        this.router.navigateByUrl(CommonRoutes.SIGN_IN);
      });
    });
  }

  canActivateChild(childRoute: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> {
    return this.canActivate(childRoute, state);
  }

  isLogin(): Observable<any> {
    return this.authService.getProfile();
  }

}
