import {Component, OnInit} from '@angular/core';
import {SettingService} from '../setting.service';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {NotificationService} from './notification.service';

@Component({
  selector: 'app-notification',
  templateUrl: './notification.component.html',
  styleUrls: ['./notification.component.css']
})
export class NotificationComponent implements OnInit {

  constructor(private settingService: SettingService, private alert: CommonAlertService, private notificationService: NotificationService) {
  }

  notifications;
  email = {};
  dingTalk = {};
  workWeixin = {};
  loading = false;
  emailValid = false;
  workWeixinInValid = true;

  ngOnInit() {
    this.listSettings('email');
    this.listSettings('dingTalk');
    this.listSettings('workWeixin');
  }


  listSettings(tab) {
    this.loading = true;
    this.settingService.getSettingsByTab(tab).subscribe(data => {
      if (tab === 'email') {
        this.email = data;
      }
      if (tab === 'dingTalk') {
        this.dingTalk = data;
      }
      if (tab === 'workWeixin') {
        this.workWeixin = data;
      }
      this.loading = false;
    });
  }

  onSubmit(data, tab) {
    this.settingService.updateSettings(data, tab).subscribe(res => {
      this.alert.showAlert('修改成功！', AlertLevels.SUCCESS);
    });
  }

  onCancel(tab) {
    this.listSettings(tab);
  }

  checkEmail() {
    this.notificationService.emailCheck(this.email).subscribe(data => {
      this.emailValid = true;
      this.alert.showAlert(data['msg'], AlertLevels.SUCCESS);
    }, error => {
      this.emailValid = false;
      this.alert.showAlert(error.error.msg, AlertLevels.ERROR);
    });
  }

  checkWorkWeixin() {
    this.notificationService.workWeixinCheck(this.workWeixin).subscribe(data => {
      this.workWeixinInValid = false;
      this.alert.showAlert(data['msg'], AlertLevels.SUCCESS);
    }, error => {
      this.workWeixinInValid = true;
      this.alert.showAlert(error.error.msg, AlertLevels.ERROR);
    });
  }
}
