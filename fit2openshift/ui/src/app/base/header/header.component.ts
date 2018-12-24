import {Component, OnInit} from '@angular/core';
import {SessionService} from '../../shared/session.service';
import {Router} from '@angular/router';
import {CommonRoutes} from '../../shared/shared.const';
import {SessionUser} from '../../shared/session-user';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit {
  user: SessionUser = new SessionUser();
  username = 'guest';

  constructor(private sessionService: SessionService, private router: Router) {
  }

  ngOnInit() {
    this.getCurrentUser();
  }

  getCurrentUser() {
    this.user = this.sessionService.getCacheUser();
    if (this.user) {
      this.username = this.user.username;
    }
  }

  logOut() {
    this.sessionService.clear();
    this.router.navigateByUrl(CommonRoutes.SIGN_IN);
  }

}
