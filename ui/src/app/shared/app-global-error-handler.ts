import {ErrorHandler, Injectable} from '@angular/core';
import {CommonAlertService} from '../layout/common-alert/common-alert.service';
import {AlertLevels} from '../layout/common-alert/alert';
import {Router} from '@angular/router';
import {CommonRoutes} from '../globals';

@Injectable()
export class AppGlobalErrorHandler implements ErrorHandler {

    constructor(private commonAlert: CommonAlertService, private router: Router,) {
    }

    handleError(error) {
        console.log(error);
        if (error.status && error.status === 401) {
            this.commonAlert.showAlert(error.statusText, AlertLevels.ERROR);
            this.router.navigateByUrl(CommonRoutes.LOGIN).then(r => console.log('logout success'));
        }
    }
}
