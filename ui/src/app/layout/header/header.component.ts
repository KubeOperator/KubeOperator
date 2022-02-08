import {Component, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {SessionService} from '../../shared/auth/session.service';
import {SessionUser} from '../../shared/auth/session-user';
import {Router} from '@angular/router';
import {CommonRoutes} from '../../constant/route';
import {PasswordComponent} from './password/password.component';
import {AboutComponent} from './about/about.component';
import {NoticeService} from '../../business/message-center/mailbox/notice.service';
import {LicenseService} from '../../business/setting/license/license.service';
import {TranslateService} from '@ngx-translate/core';
import {BusinessLicenseService} from '../../shared/service/business-license.service';

@Component({
    selector: 'app-header',
    templateUrl: './header.component.html',
    styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit, OnDestroy {

    user: SessionUser = new SessionUser();
    unreadAlert = 0;
    unreadInfo = 0;
    hasLicense = false;
    haveNotices = false;
    language: string;

    @ViewChild(PasswordComponent, {static: true})
    password: PasswordComponent;

    @ViewChild(AboutComponent, {static: true})
    about: AboutComponent;
    logo: string;
    timer;

    constructor(private sessionService: SessionService, private router: Router,
                private noticeService: NoticeService, private businessLicenseService: BusinessLicenseService,
                private translateService: TranslateService) {
    }

    ngOnInit(): void {
        this.getProfile();
        const currentLanguage = localStorage.getItem('currentLanguage');
        if (currentLanguage) {
            this.language = currentLanguage;
        } else {
            this.language = 'zh-CN';
        }
    }

    getProfile() {
        this.sessionService.getProfile().subscribe(res => {
            if (res != null) {
                this.user = res.user;
                this.hasLicense = this.businessLicenseService.licenseValid;
                if (this.hasLicense) {
                    this.listUnreadMsg(this.user.name);
                    this.timer = setInterval(() => {
                        this.listUnreadMsg(this.user.name);
                    }, 60000);
                }
            }
        }, error => {
            this.sessionService.clear();
            this.router.navigateByUrl(CommonRoutes.LOGIN).then();
        })
    }

    ngOnDestroy() {
        if (this.timer) {
            clearInterval(this.timer);
        }
    }

    changePassword() {
        this.password.open(this.user);
    }

    listUnreadMsg(userName) {
        this.noticeService.listUnread().subscribe(res => {
            this.unreadAlert = res.warning;
            this.unreadInfo = res.info;
            if (this.unreadAlert > 0 || this.unreadAlert > 0) {
                this.haveNotices = true;
            }
        }, error => {
        });
    }

    openDoc() {
        window.open('https://kubeoperator.io/docs/', '_blank');
    }

    openSwagger() {
        window.open('/swagger/index.html', '_blank');
    }

    setLogo(logo: string) {
        this.logo = logo;
    }

    logOut() {
        this.sessionService.logout().toPromise().finally(
            () => {
                this.sessionService.clear();
                this.router.navigateByUrl(CommonRoutes.LOGIN).then();
            }
        );
    }

    openAbout() {
        this.about.open();
    }

    changeLanguage(language) {
        localStorage.setItem('currentLanguage', language);
        this.translateService.use(language);
        this.language = language;
        window.location.reload();
    }
}

