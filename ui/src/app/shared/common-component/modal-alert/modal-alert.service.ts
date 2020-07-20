import {Injectable} from '@angular/core';
import {Subject} from 'rxjs';
import {Alert, AlertLevels} from '../../../layout/common-alert/alert';

@Injectable({
    providedIn: 'root'
})
export class ModalAlertService {

    constructor() {
    }

    alertQueue = new Subject<Alert>();
    $alertQueue = this.alertQueue.asObservable();

    showAlert(error: any, level: AlertLevels) {
        let msg = '';
        if ((typeof error).toLowerCase() === 'string') {
            msg = error;
        } else if (error.error.msg != null) {
            msg = error.error.msg;
        } else {
            msg = error.error;
        }
        this.alertQueue.next(new Alert(msg, level));
    }
}
