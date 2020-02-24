import {Component, OnInit} from '@angular/core';
import {SessionUser} from '../../shared/session-user';
import {SessionService} from '../../shared/session.service';

@Component({
  selector: 'app-navigator',
  templateUrl: './navigator.component.html',
  styleUrls: ['./navigator.component.css']
})
export class NavigatorComponent implements OnInit {

  user: SessionUser;

  constructor(private sessionService: SessionService) {
  }

  getCurrentUser() {
    this.user = this.sessionService.getCacheUser();
    console.log(this.user);
  }

  ngOnInit() {
    this.getCurrentUser();
  }

}
