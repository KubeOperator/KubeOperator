import {Injectable} from '@angular/core';
import {Subject} from 'rxjs';
import {Alert, AlertLevels} from './alert';

@Injectable({
    providedIn: 'root'
})
export class AppAlertService {

    constructor() {
    }

    alertQueue = new Subject<Alert>();
    $alertQueue = this.alertQueue.asObservable();


    showAlert(msg: string, level: AlertLevels) {
        this.alertQueue.next(new Alert(msg, level));
    }
}
