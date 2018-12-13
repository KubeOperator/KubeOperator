import { Component, OnInit, ViewChild } from '@angular/core';

import { HostCreateComponent } from './host-create/host-create.component';
import { HostListComponent } from './host-list/host-list.component';

@Component({
  selector: 'app-host',
  templateUrl: './host.component.html',
  styles: []
})
export class HostComponent implements OnInit {
  @ViewChild(HostCreateComponent)
  creationHost: HostCreateComponent;
  @ViewChild(HostListComponent)
  listHost: HostListComponent;

  constructor() { }

  ngOnInit() {
  }

  openModal() {
    this.creationHost.newHost();
  }

  onHostCreated() {
    this.listHost.refresh();
  }
}
