import {Component, OnInit} from '@angular/core';
import {Setting} from './setting';
import {SettingService} from './setting.service';
import {TipService} from '../tip/tip.service';
import {TipLevels} from '../tip/tipLevels';

@Component({
  selector: 'app-setting',
  templateUrl: './setting.component.html',
  styleUrls: ['./setting.component.css']
})
export class SettingComponent implements OnInit {

  orgSettings: Setting[] = [];
  settings: Setting[] = [];

  constructor(private  settingService: SettingService, private tipService: TipService) {
  }

  ngOnInit() {
    this.listSettings();
  }

  listSettings() {
    this.settingService.listSettings().subscribe(data => {
      this.settings = data;
      console.log(this.settings);
      this.orgSettings = JSON.parse(JSON.stringify(this.settings));
    });
  }


  onCancel() {
    this.listSettings();
  }

  onSubmit() {
    this.orgSettings.forEach(os => {
      this.settings.forEach(s => {
        if (os.key === s.key && os.value !== s.value) {
          this.settingService.updateSetting(s.key, s).subscribe(data => {
            this.tipService.showTip('修改成功！', TipLevels.SUCCESS);
          }, err => {
            this.tipService.showTip('修改失败！:' + err, TipLevels.ERROR);
          });
        }
      });
    });
  }
}
