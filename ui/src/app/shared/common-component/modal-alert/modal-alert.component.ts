import {Component, OnInit} from '@angular/core';
import {Alert, AlertLevels} from '../../../layout/common-alert/alert';
import {ModalAlertService} from './modal-alert.service';

@Component({
    selector: 'app-modal-alert',
    templateUrl: './modal-alert.component.html',
    styleUrls: ['./modal-alert.component.css']
})
export class ModalAlertComponent implements OnInit {

    msg = '';
    show = false;
    currentAlert: Alert;
    level = '';
    leftSeconds = 5;
    defaultLeftTime = 5;
    defaultInterval = 1000;

    constructor(private modalAlertService: ModalAlertService) {
    }

    ngOnInit(): void {
        this.showTip();
    }

    showTip() {
        this.modalAlertService.$alertQueue.subscribe(alert => {
            this.currentAlert = alert;
            switch (alert.level) {
                case AlertLevels.SUCCESS:
                    this.level = 'success';
                    this.msg = alert.msg;
                    break;
                case AlertLevels.ERROR:
                    this.level = 'danger';
                    this.msg = alert.msg;
            }
            this.show = true;
            const timer = setInterval(() => {
                this.leftSeconds--;
                if (this.leftSeconds < 0 || !this.modalAlertService) {
                    this.show = false;
                    clearInterval(timer);
                    this.leftSeconds = this.defaultLeftTime;
                }
            }, this.defaultInterval);
        });
    }
}
