import {Injectable} from '@angular/core';
import {SessionService} from './session.service';
import {
    ActivatedRouteSnapshot,
    CanActivate,
    CanActivateChild,
    Router,
    RouterStateSnapshot,
    UrlTree
} from '@angular/router';
import {ModalAlertService} from '../common-component/modal-alert/modal-alert.service';
import {Observable} from 'rxjs';
import {AlertLevels} from '../../layout/common-alert/alert';
import {CommonRoutes} from '../../constant/route';

@Injectable({
    providedIn: 'root'
})
export class AdminAuthService implements CanActivateChild, CanActivate {

    constructor(private sessionService: SessionService, private router: Router, private modalAlertService: ModalAlertService) {
    }

    canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot):
        Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
        return new Observable<boolean>((observer) => {
            this.sessionService.getProfile().subscribe(res => {
                if (res != null && res.user.isAdmin) {
                    localStorage.setItem("cs", res.cs)
                    observer.next(true);
                    observer.complete();
                } else {
                    this.modalAlertService.showAlert('no profile', AlertLevels.ERROR);
                    this.router.navigateByUrl(CommonRoutes.KO_ROOT).then();
                }
            }, error => {
                this.sessionService.clear();
                this.router.navigateByUrl(CommonRoutes.LOGIN).then();
            })
        });
    }

    canActivateChild(childRoute: ActivatedRouteSnapshot, state: RouterStateSnapshot):
        Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
        return this.canActivate(childRoute, state);
    }
}
