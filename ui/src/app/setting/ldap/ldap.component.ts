import {Component, OnInit} from '@angular/core';
import {Settings} from '../setting';
import {SettingService} from '../setting.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {CommonAlertService} from '../../base/header/common-alert.service';

@Component({
  selector: 'app-ldap',
  templateUrl: './ldap.component.html',
  styleUrls: ['./ldap.component.css']
})
export class LdapComponent implements OnInit {

  constructor(private settingService: SettingService, private alert: CommonAlertService) {
  }

  settings: Settings = new Settings();

  ngOnInit() {
    this.listSettings();
  }

  listSettings() {
    this.settingService.getSettingsByTab('ldap').subscribe(data => {
      this.settings = data;
    });
  }

  onCancel() {
    this.listSettings();
  }

  onSubmit() {
    this.settingService.updateSettings(this.settings, 'ldap').subscribe(data => {
      this.settings = data;
      this.alert.showAlert('修改成功！', AlertLevels.SUCCESS);
    });
  }


}
