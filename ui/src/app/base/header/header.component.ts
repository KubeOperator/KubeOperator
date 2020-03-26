import {Component, OnInit, ViewChild} from '@angular/core';
import {SessionService} from '../../shared/session.service';
import {Router} from '@angular/router';
import {CommonRoutes} from '../../shared/shared.const';
import {Profile, SessionUser} from '../../shared/session-user';
import {PasswordComponent} from './components/password/password.component';
import {BaseService} from '../base.service';
import {Version} from './version';
import {MessageCenterService} from '../../message-center/message-center.service';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit {
  profile: Profile = new Profile();
  username = 'guest';
  showVersionInfo = false;
  version: Version = new Version();
  @ViewChild(PasswordComponent, {static: true})
  password: PasswordComponent;
  info;
  warning;
  timer;


  constructor(private sessionService: SessionService, private router: Router, private baseService: BaseService,
              private messageCenterService: MessageCenterService) {
  }

  ngOnInit() {
    this.getProfile();
    this.getVersionInfo();
    this.getUnReadMessage();

    this.timer = setInterval(() => {
      this.getUnReadMessage();
    }, 60000);
  }

  // tslint:disable-next-line:use-lifecycle-interface
  ngOnDestroy() {
    if (this.timer) {
      clearInterval(this.timer);
    }
  }

  getProfile() {
    this.profile = this.sessionService.getCacheProfile();
    if (this.profile) {
      this.username = this.profile.user.username;
    }
  }

  getVersionInfo() {
    this.baseService.getVersion().subscribe(data => {
      this.version = data;
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

  refreshCache() {
    this.profile = this.sessionService.getCacheProfile();
  }

  getUnReadMessage() {
    this.messageCenterService.unReadMessage().subscribe(res => {
      this.info = res.info;
      this.warning = res.warning;
    });
  }

  toMessage() {
    this.router.navigateByUrl('/messageCenter/localMail');
  }
}
