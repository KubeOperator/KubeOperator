import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Zone} from '../zone';
import {ZoneService} from '../zone.service';
import * as ipaddr from 'ipaddr.js';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';

@Component({
  selector: 'app-zone-edit',
  templateUrl: './zone-edit.component.html',
  styleUrls: ['./zone-edit.component.css']
})
export class ZoneEditComponent implements OnInit {

  item: Zone;
  @Input() open = false;
  @Output() openChange = new EventEmitter();
  networkErrors = [];
  @Output() completed = new EventEmitter();

  constructor(private zoneService: ZoneService, private alert: CommonAlertService) {
  }

  ngOnInit() {
  }

  onConfirm() {
    this.zoneService.updateZones(this.item).subscribe(data => {
      this.onCancel();
      this.alert.showAlert('编辑可用区成功', AlertLevels.SUCCESS);
      this.completed.emit();
    });
  }

  checkParams() {
    const start_ip = this.item.vars['ip_start'];
    const end_ip = this.item.vars['ip_end'];
    let result = true;
    this.networkErrors = [];

    if (!ipaddr.isValid(start_ip)) {
      result = false;
      this.networkErrors.push('起始IP不是有效的IPV4地址!');
    }
    if (!ipaddr.isValid(end_ip)) {
      result = false;
      this.networkErrors.push('截止IP不是有效的IPV4地址!');
    }
    if (ipaddr.isValid(start_ip) && ipaddr.isValid(end_ip)) {
      const start = ipaddr.parse(start_ip);
      const end = ipaddr.parse(end_ip);
      if (start > end) {
        result = false;
        this.networkErrors.push('截止IP必须大于起始IP!');
      }
    }
    return result;
  }

  onCancel() {
    this.open = false;
    this.openChange.emit(this.open);
  }

}
