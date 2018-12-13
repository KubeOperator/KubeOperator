import { Injectable } from '@angular/core';
import {
  CanActivate, Router,
  ActivatedRouteSnapshot,
  RouterStateSnapshot,
  CanActivateChild
} from '@angular/router';
import { Observable } from "rxjs/index";
import { AuthService } from "../auth.service";


@Injectable({
  providedIn: 'root'
})
export class AuthCheckGuard implements CanActivate, CanActivateChild {
  constructor(private router: Router, private authService: AuthService) { }

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot,
  ): Observable<boolean>|boolean {
    // return of(true);
    const originPath = window.location.pathname;
    const redirectUrl = '/admin/login/?next=' + originPath;
    return new Observable<boolean>((observer) => {
      const user = this.authService.getCurrentUser();
      if (!user){
        this.authService.getProfile()
          .subscribe(
            user => {
              observer.next(true);
              observer.complete();
            },
            () => {
              observer.next(false);
              observer.complete();
              window.location.href = redirectUrl;
            }
          )
      }
      else {
        observer.next(true);
        observer.complete();
      }
      return {unsubscribe(){}}
    })
  }

  canActivateChild(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot):  Observable<boolean>|boolean {
    return this.canActivate(route, state)
  }
}
