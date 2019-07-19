import {Component, OnInit, ViewChild} from '@angular/core';
import {SessionService} from '../../shared/session.service';
import {Router} from '@angular/router';
import {CommonRoutes} from '../../shared/shared.const';
import {SessionUser} from '../../shared/session-user';
import {PasswordComponent} from './components/password/password.component';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit {
  user: SessionUser = new SessionUser();
  username = 'guest';

  @ViewChild(PasswordComponent)
  password: PasswordComponent;

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

  changePassword() {
    this.password.opened = true;
  }

}
