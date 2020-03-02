import {Component, OnInit} from '@angular/core';
import {Profile, SessionUser} from '../../shared/session-user';
import {SessionService} from '../../shared/session.service';
import {User} from '../../user/user';

@Component({
  selector: 'app-navigator',
  templateUrl: './navigator.component.html',
  styleUrls: ['./navigator.component.css']
})
export class NavigatorComponent implements OnInit {

  user: SessionUser;
  showItem = false;

  constructor(private sessionService: SessionService) {
  }

  getProfile() {
    this.user = this.sessionService.getCacheProfile().user;

    this.sessionService.getProfile().subscribe(data => {
      const profile = data;
      this.user = profile.user;
      if (this.user.is_superuser) {
        this.showItem = true;
      } else {
        for (const item of profile.items) {
          for (const rm of profile.item_role_mappings) {
            if (item.name === rm.item_name && rm.role !== 'VIEWER') {
              this.showItem = true;
              break;
            }
          }
        }
      }
    });
  }

  ngOnInit() {
    this.getProfile();
  }

}
