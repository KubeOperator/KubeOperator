import { Component, OnInit, ViewChild } from '@angular/core';

import { GroupCreateComponent } from './group-create/group-create.component';
import { GroupListComponent } from './group-list/group-list.component';

@Component({
  selector: 'app-group',
  templateUrl: './group.component.html',
  styles: []
})
export class GroupComponent implements OnInit {
  @ViewChild(GroupCreateComponent)
  creationGroup: GroupCreateComponent;
  @ViewChild(GroupListComponent)
  listGroup: GroupListComponent;

  constructor() { }

  ngOnInit() {
  }

  newGroup() {
    this.creationGroup.newGroup();
  }

  onGroupCreated() {
    this.listGroup.refresh();
  }

}
