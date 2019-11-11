import {Component, OnInit} from '@angular/core';
import {Settings} from '../setting';
import {SettingService} from '../setting.service';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';

@Component({
  selector: 'app-system-setting',
  templateUrl: './system-setting.component.html',
  styleUrls: ['./system-setting.component.css']
})
export class SystemSettingComponent implements OnInit {


  constructor(private  settingService: SettingService, private alert: CommonAlertService) {
  }

  settings: Settings;

  ngOnInit() {
    this.listSettings();
  }

  listSettings() {
    this.settingService.getSettings().subscribe(data => {
      this.settings = data;
    });
  }


  onCancel() {
    this.listSettings();
  }

  onSubmit() {
    this.settingService.updateSettings(this.settings).subscribe(data => {
      this.settings = data;
      this.alert.showAlert('修改成功！', AlertLevels.SUCCESS);
    });
  }
}
