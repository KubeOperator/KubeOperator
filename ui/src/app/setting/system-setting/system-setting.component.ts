import {Component, OnInit} from '@angular/core';
import {Setting} from '../setting';
import {SettingService} from '../setting.service';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';

@Component({
  selector: 'app-system-setting',
  templateUrl: './system-setting.component.html',
  styleUrls: ['./system-setting.component.css']
})
export class SystemSettingComponent implements OnInit {

  orgSettings: Setting[] = [];
  settings: Setting[] = [];

  constructor(private  settingService: SettingService, private alert: CommonAlertService) {
  }

  ngOnInit() {
    this.listSettings();
  }

  listSettings() {
    this.settingService.listSettings().subscribe(data => {
      this.settings = data;
      this.orgSettings = JSON.parse(JSON.stringify(this.settings));
    });
  }


  onCancel() {
    this.listSettings();
  }

  onSubmit() {
    this.orgSettings.forEach(os => {
      this.settings.forEach(s => {
        if (os.key === s.key && os.value !== s.value && this.validate(s) ) {
          this.settingService.updateSetting(s.key, s).subscribe(data => {
            this.alert.showAlert('修改成功！', AlertLevels.SUCCESS);
          }, err => {
            this.alert.showAlert('修改失败！:' + err, AlertLevels.ERROR);
          });
        }
      });
    });
  }

  validate(setting) {
    const ipReg =  /^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$/g;
    if (setting.key === 'local_hostname') {
      const validate: boolean = ipReg.test(setting.value);
      if (!validate) {
        this.alert.showAlert('请输入正确的IP地址！', AlertLevels.ERROR);
        return false;
      }
    }
    const domainReg = /(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]/g;
    if (setting.key === 'domain_suffix') {
      const validate: boolean = domainReg.test(setting.value);
      console.log(validate);
      console.log(setting.value);
      if (!validate) {
        this.alert.showAlert('请输入正确的域名后缀！', AlertLevels.ERROR);
        return false;
      }
    }
    return true;
  }

}
