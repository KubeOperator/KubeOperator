import {Component, OnInit} from '@angular/core';
import {CommonAlertService} from '../../common-alert.service';
import {Alert, AlertLevels} from './alert';

const defaultInterval = 1000;
const defaultLeftTime = 5;

@Component({
  selector: 'app-common-alert',
  templateUrl: './common-alert.component.html',
  styleUrls: ['./common-alert.component.css']
})
export class CommonAlertComponent implements OnInit {

  constructor(private commonAlertService: CommonAlertService) {
  }

  show = false;
  level = '';
  msg = '';
  leftSeconds: number = defaultLeftTime;

  currentAlert: Alert;


  ngOnInit() {
    this.showTip();
  }

  showTip() {
    this.commonAlertService.$alertQueue.subscribe(alert => {
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
        if (this.leftSeconds < 0 || !this.commonAlertService) {
          this.show = false;
          clearInterval(timer);
          this.leftSeconds = defaultLeftTime;
        }
      }, defaultInterval);
    });
  }

  close() {
    this.show = false;
  }
}
