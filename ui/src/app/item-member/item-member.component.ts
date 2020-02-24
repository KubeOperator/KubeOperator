import {Component, OnInit, ViewChild} from '@angular/core';
import {ItemMemberListComponent} from './item-member-list/item-member-list.component';
import {ItemMemberCreateComponent} from './item-member-create/item-member-create.component';
import {ActivatedRoute} from '@angular/router';
import {Profile} from '../shared/session-user';

@Component({
  selector: 'app-item-member',
  templateUrl: './item-member.component.html',
  styleUrls: ['./item-member.component.css']
})
export class ItemMemberComponent implements OnInit {
  currentItemName: string;

  constructor(private route: ActivatedRoute) {
  }

  @ViewChild(ItemMemberListComponent, {static: true})
  list: ItemMemberListComponent;

  @ViewChild(ItemMemberCreateComponent, {static: true})
  creation: ItemMemberCreateComponent;


  ngOnInit() {
    this.currentItemName = this.route.snapshot.queryParams['name'];
  }

  openAdd(profiles: Profile[]) {
    this.creation.open(profiles);
  }

  postAdd() {
    this.list.refresh();
  }

}
