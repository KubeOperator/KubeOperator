import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import * as globals from '../../globals';
import {Ceph} from '../ceph';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {CephService} from '../ceph.service';
import {NgForm} from '@angular/forms';


@Component({
  selector: 'app-ceph-create',
  templateUrl: './ceph-create.component.html',
  styleUrls: ['./ceph-create.component.css']
})
export class CephCreateComponent implements OnInit {

  opened = false;
  name_pattern = globals.host_name_pattern;
  name_pattern_tip = globals.host_name_pattern_tip;
  item: Ceph = new Ceph();
  isSubmitGoing = false;
  @Output() create = new EventEmitter<boolean>();
  @ViewChild('itemForm', {static: true}) itemFrom: NgForm;

  constructor(private cephService: CephService, private alertService: CommonAlertService) {
  }

  ngOnInit() {
  }

  open() {
    this.opened = true;
    this.item = new Ceph();
    this.item.vars['ceph_imageFormat'] = '2';
    this.item.vars['ceph_fsType'] = 'ext4';
  }

  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.isSubmitGoing = true;
    this.cephService.create(this.item).subscribe(data => {
      this.isSubmitGoing = false;
      this.opened = false;
      this.create.emit(true);
      this.alertService.showAlert('创建 Ceph 成功', AlertLevels.SUCCESS);
    }, error1 => {
      this.isSubmitGoing = false;
      this.alertService.showAlert('创建 Ceph 失败', AlertLevels.ERROR);
    });
  }

  onCancel() {
    this.opened = false;
  }
}
