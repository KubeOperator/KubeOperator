import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {NgForm} from '@angular/forms';
import {NfsStorage} from '../nfs';
import {Host} from '../../host/host';
import {NfsService} from '../nfs.service';
import {HostService} from '../../host/host.service';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import * as globals from '../../globals';


@Component({
  selector: 'app-nfs-create',
  templateUrl: './nfs-create.component.html',
  styleUrls: ['./nfs-create.component.css']
})
export class NfsCreateComponent implements OnInit {

  constructor(private nfsService: NfsService, private hostService: HostService, private alertService: CommonAlertService) {
  }

  opened = false;
  isSubmitGoing = false;
  item: NfsStorage = new NfsStorage();
  hosts: Host[] = [];
  @Output() create = new EventEmitter<boolean>();
  @ViewChild('itemForm', {static: true}) itemFrom: NgForm;
  name_pattern = globals.host_name_pattern;
  name_pattern_tip = globals.host_name_pattern_tip;


  ngOnInit() {
  }

  onCancel() {
    this.opened = false;
  }

  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.isSubmitGoing = true;
    this.nfsService.create(this.item).subscribe(data => {
      this.isSubmitGoing = false;
      this.opened = false;
      this.create.emit(true);
      this.alertService.showAlert('创建 NFS 成功', AlertLevels.SUCCESS);
    });
  }

  listHosts() {
    this.hostService.listAll().subscribe(data => {
      this.hosts = data.filter((host) => host.cluster === null);
    });
  }

  newItem() {
    this.listHosts();
    this.itemFrom.resetForm();
    this.item = new NfsStorage();
    this.opened = true;
  }

}
