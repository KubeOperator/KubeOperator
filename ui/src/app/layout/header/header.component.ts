import {Component, OnInit, ViewChild} from '@angular/core';
import {SessionService} from '../../shared/auth/session.service';
import {SessionUser} from '../../shared/auth/session-user';
import {Router} from '@angular/router';
import {CommonRoutes} from '../../constant/route';
import {PasswordComponent} from './password/password.component';
import {AboutComponent} from './about/about.component';
import {NoticeService} from '../../business/message-center/mailbox/notice.service';
import {LicenseService} from '../../business/setting/license/license.service';

@Component({
    selector: 'app-header',
    templateUrl: './header.component.html',
    styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit {

    user: SessionUser = new SessionUser();
    unreadAlert = 0;
    unreadInfo = 0;
    hasLicense = false;
    haveNotices = false;

    @ViewChild(PasswordComponent, {static: true})
    password: PasswordComponent;

    @ViewChild(AboutComponent, {static: true})
    about: AboutComponent;
    logo: string;
    timer;

    constructor(private sessionService: SessionService, private router: Router,
                private noticeService: NoticeService, private licenseService: LicenseService) {
    }

    ngOnInit(): void {
        this.getProfile();

    }

    getProfile() {
        const profile = this.sessionService.getCacheProfile();
        if (profile != null) {
            this.user = profile.user;
            this.licenseService.get().subscribe(data => {
                this.hasLicense = true;
                this.listUnreadMsg(this.user.name);
                this.timer = setInterval(() => {
                    this.listUnreadMsg(this.user.name);
                }, 60000);
            });
        }
    }

    // tslint:disable-next-line:use-lifecycle-interface
    ngOnDestroy() {
        if (this.timer) {
            clearInterval(this.timer);
        }
    }

    changePassword() {
        this.password.open(this.user);
    }

    listUnreadMsg(userName) {
        this.noticeService.listUnread(userName).subscribe(res => {
            this.unreadAlert = res.warning;
            this.unreadInfo = res.info;
            if (this.unreadAlert > 0 || this.unreadAlert > 0) {
                this.haveNotices = true;
            }
        }, error => {
        });
    }

    openDoc() {
        window.open('https://kubeoperator.io/docs/', 'blank');
    }

    openSwagger() {
        window.open('/swagger/index.html', 'blank');
    }

    setLogo(logo: string) {
        this.logo = logo;
    }

    logOut() {
        this.sessionService.clear();
        this.router.navigateByUrl(CommonRoutes.LOGIN).then();
    }

    openAbout() {
        this.about.open();
    }
}
