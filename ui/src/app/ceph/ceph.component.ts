import {Component, OnInit, ViewChild} from '@angular/core';
import {CephListComponent} from "./ceph-list/ceph-list.component";
import {CephCreateComponent} from "./ceph-create/ceph-create.component";

@Component({
  selector: 'app-ceph',
  templateUrl: './ceph.component.html',
  styleUrls: ['./ceph.component.css']
})
export class CephComponent implements OnInit {

  @ViewChild(CephListComponent, {static: true})
  list: CephListComponent;

  @ViewChild(CephCreateComponent, {static: true})
  creation: CephCreateComponent;

  constructor() {
  }

  ngOnInit() {
  }


  openModal() {
    this.creation.open();
  }

  create(created: boolean) {
    if (created) {
      this.refresh();
    }
  }

  refresh() {
    this.list.refresh();
  }
}
