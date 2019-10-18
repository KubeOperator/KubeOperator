import {Component, OnInit, ViewChild} from '@angular/core';
import {NfsListComponent} from './nfs-list/nfs-list.component';
import {NfsCreateComponent} from './nfs-create/nfs-create.component';

@Component({
  selector: 'app-nfs',
  templateUrl: './nfs.component.html',
  styleUrls: ['./nfs.component.css']
})
export class NfsComponent implements OnInit {

  @ViewChild(NfsListComponent, {static: true})
  list: NfsListComponent;

  @ViewChild(NfsCreateComponent, {static: true})
  creation: NfsCreateComponent;

  constructor() {
  }

  ngOnInit() {
  }

  openModal() {
    this.creation.newItem();
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
