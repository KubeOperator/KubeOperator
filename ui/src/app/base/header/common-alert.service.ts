import {Injectable} from '@angular/core';
import {Subject} from 'rxjs';
import {Alert, AlertLevels} from './components/common-alert/alert';

@Injectable({
  providedIn: 'root'
})
export class CommonAlertService {

  alertQueue = new Subject<Alert>();
  $alertQueue = this.alertQueue.asObservable();

  constructor() {
  }

  showAlert(msg: string, level: AlertLevels) {
    this.alertQueue.next(new Alert(msg, level));
  }
}
