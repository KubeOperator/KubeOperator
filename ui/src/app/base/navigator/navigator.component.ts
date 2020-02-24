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

  constructor(private sessionService: SessionService) {
  }

  getProfile() {
    const profile = this.sessionService.getCacheProfile();
    this.user = profile.user;
    console.log(this.user);
  }

  ngOnInit() {
    this.getProfile();
  }

}
