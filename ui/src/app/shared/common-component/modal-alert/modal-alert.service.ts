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

    showAlert(msg: string, level: AlertLevels) {
        this.alertQueue.next(new Alert(msg, level));
    }
}
