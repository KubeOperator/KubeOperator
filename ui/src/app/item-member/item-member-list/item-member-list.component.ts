import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {ItemMemberService} from '../item-member.service';
import {ItemMember} from '../item-member';
import {Profile} from '../../shared/session-user';
import {ActivatedRoute} from '@angular/router';

@Component({
  selector: 'app-item-member-list',
  templateUrl: './item-member-list.component.html',
  styleUrls: ['./item-member-list.component.css']
})
export class ItemMemberListComponent implements OnInit {

  constructor(private itemMemberService: ItemMemberService, private route: ActivatedRoute) {
  }

  itemMember: ItemMember = new ItemMember();
  loading = true;
  @Output() add = new EventEmitter<Profile[]>();
  currentItem;

  ngOnInit() {
    this.currentItem = this.route.snapshot.queryParams['name'];
    this.refresh();
  }

  onAdd() {
    this.add.emit(this.itemMember.profiles);
  }

  refresh() {
    this.loading = true;
    this.itemMemberService.getItemProfiles(this.currentItem).subscribe(data => {
      this.loading = false;
      this.itemMember = data;
    });
  }

  formatRole(p: Profile) {
    return p.item_role_mappings.find(mp => {
      return mp.item_name === this.currentItem;
    })['role'];
  }
}
