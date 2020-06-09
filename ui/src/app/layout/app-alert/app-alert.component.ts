import {Component, OnInit} from '@angular/core';
import {Alert, AlertLevels} from '../common-alert/alert';
import {AppAlertService} from './app-alert.service';

@Component({
    selector: 'app-app-alert',
    templateUrl: './app-alert.component.html',
    styleUrls: ['./app-alert.component.css']
})
export class AppAlertComponent implements OnInit {

    show = false;
    level = '';
    msg = '';
    leftSeconds = 5;
    defaultLeftTime = 5;
    defaultInterval = 1000;
    currentAlert: Alert;

    constructor(private appAlertService: AppAlertService) {
    }

    ngOnInit(): void {
        this.showTip();
    }

    showTip() {
        this.appAlertService.$alertQueue.subscribe(alert => {
            this.currentAlert = alert;
            switch (alert.level) {
                case AlertLevels.SUCCESS:
                    this.level = 'info';
                    this.msg = alert.msg;
                    break;
                case AlertLevels.ERROR:
                    this.level = 'danger';
                    this.msg = alert.msg;
            }
            this.show = true;
            const timer = setInterval(() => {
                this.leftSeconds--;
                if (this.leftSeconds < 0 || !this.appAlertService) {
                    this.show = false;
                    clearInterval(timer);
                    this.leftSeconds = this.defaultLeftTime;
                }
            }, this.defaultInterval);
        });
    }

    close() {
        this.show = false;
    }

}
