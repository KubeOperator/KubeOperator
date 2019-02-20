import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {HostService} from '../host.service';
import {Host} from '../host';
import {TipService} from '../../tip/tip.service';
import {TipLevels} from '../../tip/tipLevels';

@Component({
  selector: 'app-host-list',
  templateUrl: './host-list.component.html',
  styleUrls: ['./host-list.component.css']
})
export class HostListComponent implements OnInit {

  hosts: Host[] = [];
  loading = false;
  deleteModal = false;
  selectedHost: Host = new Host();
  @Output() addHost = new EventEmitter();


  constructor(private hostService: HostService, private tipService: TipService) {
  }

  ngOnInit() {
    this.listHost();
  }

  deleteHost(host: Host) {
    this.deleteModal = true;
    this.selectedHost = host;
  }

  confirmDelete() {
    this.hostService.deleteHost(this.selectedHost.id).subscribe(data => {
      this.deleteModal = false;
      this.refresh();
      this.tipService.showTip('删除主机成功!', TipLevels.SUCCESS);
    }, err => {
      this.tipService.showTip('删除失败:' + err, TipLevels.ERROR);
    });
  }


  refresh() {
    this.listHost();
  }

  addNewHost() {
    this.addHost.emit();
  }


  listHost() {
    this.hostService.listHosts().subscribe(data => {
      this.hosts = data;
      this.loading = false;
    }, error => {
      this.loading = false;
    });
  }
}
