import {ErrorHandler, Injectable} from '@angular/core';
import {CommonAlertService} from '../../layout/common-alert/common-alert.service';
import {AlertLevels} from '../../layout/common-alert/alert';
import {Router} from '@angular/router';
import {CommonRoutes} from '../../constant/route';
import {AppAlertService} from '../../layout/app-alert/app-alert.service';

@Injectable()
export class AppGlobalErrorHandler implements ErrorHandler {

    constructor(private appAlertService: AppAlertService,
                private commonAlertService: CommonAlertService, private router: Router) {
    }

    handleError(error) {
        switch (error.status) {
            case  400:
                this.appAlertService.showAlert(error.statusText, AlertLevels.ERROR);
                break;
            case 500:
                this.commonAlertService.showAlert(error.statusText, AlertLevels.ERROR);
                break;
            case 401:
                this.router.navigateByUrl('login').then();
        }
    }
}
