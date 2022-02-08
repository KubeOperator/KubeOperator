import {Injectable} from '@angular/core';
import {CommonRoutes} from '../../constant/route';
import {
    ActivatedRouteSnapshot,
    CanActivate,
    CanActivateChild,
    Router,
    RouterStateSnapshot,
    UrlTree
} from '@angular/router';
import {Observable} from 'rxjs';
import {SessionService} from './session.service';
import {AlertLevels} from '../../layout/common-alert/alert';
import {ModalAlertService} from '../common-component/modal-alert/modal-alert.service';

@Injectable({
    providedIn: 'root'
})
export class AuthUserService implements CanActivate, CanActivateChild {

    constructor(private sessionService: SessionService, private router: Router, private modalAlertService: ModalAlertService) {
    }

    canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot):
        Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
        return new Observable<boolean>((observer) => {
            this.sessionService.getProfile().subscribe(res => {
                if (res != null) {
                    observer.next(true);
                    observer.complete();
                } else {
                    this.modalAlertService.showAlert('no profile', AlertLevels.ERROR);
                    this.router.navigateByUrl(CommonRoutes.LOGIN).then();
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
