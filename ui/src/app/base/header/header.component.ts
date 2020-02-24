import {Component, OnInit, ViewChild} from '@angular/core';
import {SessionService} from '../../shared/session.service';
import {Router} from '@angular/router';
import {CommonRoutes} from '../../shared/shared.const';
import {Profile, SessionUser} from '../../shared/session-user';
import {PasswordComponent} from './components/password/password.component';
import {BaseService} from '../base.service';
import {Version} from './version';
import {ItemChangeComponent} from './components/item-change/item-change.component';

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

  @ViewChild(ItemChangeComponent, {static: true})
  itemChange: ItemChangeComponent;

  constructor(private sessionService: SessionService, private router: Router, private baseService: BaseService) {
  }

  ngOnInit() {
    this.getProfile();
    this.getVersionInfo();
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

  openChangeItem() {
    this.itemChange.open();
  }

  refreshCache() {
    this.profile = this.sessionService.getCacheProfile();
  }
}
