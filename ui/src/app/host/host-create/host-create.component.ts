import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Host} from '../host';
import {HostService} from '../host.service';
import {Credential} from '../../credential/credential-list/credential';
import {CredentialService} from '../../credential/credential.service';
import {NgForm} from '@angular/forms';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import * as globals from '../../globals';
import {SettingService} from '../../setting/setting.service';


@Component({
  selector: 'app-host-create',
  templateUrl: './host-create.component.html',
  styleUrls: ['./host-create.component.css']
})
export class HostCreateComponent implements OnInit {

  constructor(private hostService: HostService, private alert: CommonAlertService, private credentialService: CredentialService,
              private  settingService: SettingService) {
  }

  @Output() create = new EventEmitter<boolean>();
  staticBackdrop = true;
  closable = false;
  createHostOpened: boolean;
  isSubmitGoing = false;
  host: Host = new Host();
  loading = false;
  credentials: Credential[] = [];
  @ViewChild('hostForm', {static: true}) hostFrom: NgForm;
  name_pattern = globals.host_name_pattern;
  name_pattern_tip = globals.host_name_pattern_tip;
  localIp = '0.0.0.0';

  ngOnInit() {

  }

  getLocalIp() {
    this.settingService.getSettings().subscribe(data => {
      const hostName = data['local_hostname'];
      if (hostName) {
        this.localIp = hostName;
      }
    });
  }

  listCredential() {
    this.credentialService.listCredential().subscribe(data => {
      this.credentials = data;
    });
  }

  reset() {
    this.hostFrom.resetForm({port: 22});
    this.listCredential();
    this.getLocalIp();
  }


  onCancel() {
    this.createHostOpened = false;
  }

  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.isSubmitGoing = true;
    this.loading = true;
    this.hostService.create(this.host).subscribe(data => {
      this.createHostOpened = false;
      this.isSubmitGoing = false;
      this.create.emit(true);
      this.loading = false;
      this.alert.showAlert('创建主机成功', AlertLevels.SUCCESS);
    }, res => {
      let msg = '';
      const err = res.error;
      for (const key in err) {
        if (key) {
          msg += err[key].join('');
        }

      }
      this.createHostOpened = false;
      this.isSubmitGoing = false;
      this.create.emit(true);
      this.loading = false;
      this.alert.showAlert('创建主机失败:' + msg, AlertLevels.ERROR);
    });
  }

  newHost() {
    this.reset();
    this.host = new Host();
    this.host.port = 22;
    this.createHostOpened = true;
  }
}
