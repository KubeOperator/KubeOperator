import {Component, OnInit, ViewChild} from '@angular/core';
import {HostCreateComponent} from './host-create/host-create.component';
import {HostListComponent} from './host-list/host-list.component';
import {HostImportComponent} from './host-import/host-import.component';

@Component({
  selector: 'app-host',
  templateUrl: './host.component.html',
  styleUrls: ['./host.component.css']
})
export class HostComponent implements OnInit {

  @ViewChild(HostCreateComponent, {static: true})
  creationHost: HostCreateComponent;

  @ViewChild(HostListComponent, {static: true})
  listHost: HostListComponent;

  @ViewChild(HostImportComponent, {static: true})
  importHost: HostImportComponent;

  constructor() {
  }

  ngOnInit() {
  }

  openCreate() {
    this.creationHost.newHost();
  }

  openImport() {
    this.importHost.open();
  }

  refresh() {
    this.listHost.refresh();
  }


}
