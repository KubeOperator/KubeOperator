import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {HostService} from '../host.service';
import {Host} from '../host';
import {HostInfoComponent} from '../host-info/host-info.component';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';

@Component({
  selector: 'app-host-list',
  templateUrl: './host-list.component.html',
  styleUrls: ['./host-list.component.css']
})
export class HostListComponent implements OnInit {

  hosts: Host[] = [];
  loading = false;
  deleteModal = false;
  selectedHosts: Host[] = [];
  showHostInfo = false;
  @Output() addHost = new EventEmitter();
  @ViewChild(HostInfoComponent, {static: true})
  child: HostInfoComponent;


  constructor(private hostService: HostService, private alertService: CommonAlertService) {
  }

  ngOnInit() {
    this.listHost();
  }

  onDeleted() {
    this.deleteModal = true;
  }

  canSelectedHostsDelete(): boolean {
    if (this.selectedHosts.length === 0) {
      return false;
    }
    let result = true;
    this.selectedHosts.forEach(host => {
      if (host.cluster !== null) {
        result = false;
      }
    });
    return result;
  }

  confirmDelete() {
    const promises: Promise<{}>[] = [];
    this.selectedHosts.forEach(host => {
      promises.push(this.hostService.deleteHost(host.id).toPromise());
    });
    Promise.all(promises).then(() => {
      this.refresh();
      this.alertService.showAlert('删除主机成功！', AlertLevels.SUCCESS);
    }, (error) => {
      this.alertService.showAlert('删除主机失败:' + error, AlertLevels.ERROR);
    }).finally(() => {
      this.deleteModal = false;
      this.selectedHosts = [];
    });
  }


  refresh() {
    this.listHost();
  }

  addNewHost() {
    this.addHost.emit();
  }

  openInfo(host: Host) {
    this.showHostInfo = true;
    this.child.host = host;
  }

  listHost() {
    this.loading = true;
    this.hostService.listHosts().subscribe(data => {
      this.hosts = data;
      this.loading = false;
    }, error => {
      this.loading = false;
    });
  }

  getValueOrNone(value) {
    return value == null ? '无' : value;
  }
}
