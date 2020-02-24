import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {Item} from '../item';
import {SessionService} from '../../shared/session.service';

@Component({
  selector: 'app-item-detail',
  templateUrl: './item-detail.component.html',
  styleUrls: ['./item-detail.component.css']
})
export class ItemDetailComponent implements OnInit {

  currentItem: Item = new Item();
  permission;
  profile;

  constructor(private router: Router, private route: ActivatedRoute, private sessionService: SessionService) {
  }

  ngOnInit() {
    this.route.data.subscribe(data => {
      this.currentItem = data['item'];
      this.getProfile();
    });
  }


  backToItem() {
    this.router.navigate(['item']);
  }

  getProfile() {
    this.permission = this.sessionService.getItemPermission(this.currentItem.name);
  }
}
