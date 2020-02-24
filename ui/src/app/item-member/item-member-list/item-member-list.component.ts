import {Component, OnInit} from '@angular/core';
import {ItemMemberService} from '../item-member.service';
import {ItemMember} from '../item-member';
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
  currentItem = this.route.snapshot.queryParams['name'];


  ngOnInit() {
    this.refresh();
  }

  refresh() {
    this.itemMemberService.getItemUsers(this.currentItem).forEach(data => {
      this.itemMember = data;
    });
  }
}
