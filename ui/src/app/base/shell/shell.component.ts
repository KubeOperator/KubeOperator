import {Component, OnInit} from '@angular/core';
import {SettingService} from '../../setting/setting.service';
import {Router} from '@angular/router';

@Component({
  selector: 'app-shell',
  templateUrl: './shell.component.html',
  styleUrls: ['./shell.component.css']
})
export class ShellComponent implements OnInit {

  showAlert: boolean;

  constructor(private  settingService: SettingService, private router: Router) {
  }

  ngOnInit() {
    this.showAlert = false;
    this.settingService.getSettings().subscribe(data => {
      const hostName = data['local_hostname'];
      if (!hostName) {
        this.showAlert = true;
      }
    });
  }

  closeAlert() {
    this.showAlert = false;
  }

  toSetting() {
    const linkUrl = ['kubeOperator', 'setting', 'system'];
    this.router.navigate(linkUrl);
  }
}
