import {Component, OnInit, ViewChild} from '@angular/core';
import {SessionService} from '../../shared/session.service';
import {Router} from '@angular/router';
import {CommonRoutes} from '../../shared/shared.const';
import {SessionUser} from '../../shared/session-user';
import {PasswordComponent} from './components/password/password.component';
import {BaseService} from '../base.service';
import {Version} from './version';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit {
  user: SessionUser = new SessionUser();
  username = 'guest';
  showVersionInfo = false;
  version: Version = new Version();
  @ViewChild(PasswordComponent)
  password: PasswordComponent;

  constructor(private sessionService: SessionService, private router: Router, private baseService: BaseService) {
  }

  ngOnInit() {
    this.getCurrentUser();
    this.getVersionInfo();
  }

  getCurrentUser() {
    this.user = this.sessionService.getCacheUser();
    if (this.user) {
      this.username = this.user.username;
    }
  }

  getVersionInfo() {
    this.baseService.getVersion().subscribe(data => {
      this.version = data;
      console.log(this.version);
    });
  }

  logOut() {
    this.sessionService.clear();
    this.router.navigateByUrl(CommonRoutes.SIGN_IN);
  }

  changePassword() {
    this.password.opened = true;
  }

  showInfo() {
    this.showVersionInfo = true;
  }

}
