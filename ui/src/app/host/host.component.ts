import {Component, OnInit, ViewChild} from '@angular/core';
import {NodeCreateComponent} from '../node/node-create/node-create.component';
import {NodeListComponent} from '../node/node-list/node-list.component';
import {HostCreateComponent} from './host-create/host-create.component';
import {HostListComponent} from './host-list/host-list.component';

@Component({
  selector: 'app-host',
  templateUrl: './host.component.html',
  styleUrls: ['./host.component.css']
})
export class HostComponent implements OnInit {

  @ViewChild(HostCreateComponent, { static: true })
  creationHost: HostCreateComponent;

  @ViewChild(HostListComponent, { static: true })
  listHost: HostListComponent;

  constructor() {
  }

  ngOnInit() {
  }

  openModal() {
    this.creationHost.newHost();
  }

  createHost(created: boolean) {
    if (created) {
      this.refresh();
    }
  }

  refresh() {
    this.listHost.refresh();
  }


}
