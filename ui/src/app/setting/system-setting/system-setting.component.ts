import {Component, OnInit} from '@angular/core';
import {Settings} from '../setting';
import {SettingService} from '../setting.service';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import * as globals from '../../globals';

@Component({
  selector: 'app-system-setting',
  templateUrl: './system-setting.component.html',
  styleUrls: ['./system-setting.component.css']
})
export class SystemSettingComponent implements OnInit {


  domain_pattern = globals.domain_pattern;

  constructor(private  settingService: SettingService, private alert: CommonAlertService) {
  }

  settings: Settings;

  ngOnInit() {
    this.listSettings();
  }

  listSettings() {
    this.settingService.getSettingsByTab('system').subscribe(data => {
      this.settings = data;
    });
  }


  onCancel() {
    this.listSettings();
  }

  onSubmit() {
    if (!this.validate(this.settings)) {
      return;
    }
    this.settingService.updateSettings(this.settings, 'system').subscribe(data => {
      this.settings = data;
      this.alert.showAlert('修改成功！', AlertLevels.SUCCESS);
    });
  }

  validate(setting) {
    const ipReg = /^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/;
    if (setting.local_hostname !== undefined) {
      const validate: boolean = ipReg.test(setting.local_hostname);
      if (!validate) {
        this.alert.showAlert('请输入正确的IP地址！', AlertLevels.ERROR);
        return false;
      }
    }
    const domainReg = /(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]/g;
    if (setting.domain_suffix !== undefined) {
      const validate: boolean = domainReg.test(setting.domain_suffix);
      if (!validate) {
        this.alert.showAlert('请输入正确的域名后缀！', AlertLevels.ERROR);
        return false;
      }
    }
    return true;
  }
}
