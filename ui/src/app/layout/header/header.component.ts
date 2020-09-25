import {Component, OnInit, ViewChild} from '@angular/core';
import {SessionService} from '../../shared/auth/session.service';
import {SessionUser} from '../../shared/auth/session-user';
import {Router} from '@angular/router';
import {CommonRoutes} from '../../constant/route';
import {PasswordComponent} from './password/password.component';
import {AboutComponent} from './about/about.component';
import {Theme} from '../../business/setting/theme/theme';
import { NoticeService } from '../../business/message-center/mailbox/notice.service';

@Component({
    selector: 'app-header',
    templateUrl: './header.component.html',
    styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit {

    user: SessionUser = new SessionUser();
    haveNotices: boolean; // testing purpose
    unreadAlert: number;
    unreadInfo: number;

    @ViewChild(PasswordComponent, {static: true})
    password: PasswordComponent;

    @ViewChild(AboutComponent, {static: true})
    about: AboutComponent;
    logo: string;

    constructor(private sessionService: SessionService, private router: Router,
                private noticeService: NoticeService) {
    }

    ngOnInit(): void {
        this.getProfile();
        this.noticeService.currentUnread.subscribe(data => {
            this.unreadAlert = data.unreadAlert;
            this.unreadInfo = data.unreadInfo;
            this.checkNotice();
        });
    }

    // Testing purpose
    checkNotice() {
        if (this.unreadAlert + this.unreadInfo > 0) {
            this.haveNotices = true;
        } else {
            this.haveNotices = false;
        }
    }
    //

    getProfile() {
        const profile = this.sessionService.getCacheProfile();
        if (profile != null) {
            this.user = profile.user;
        }
    }

    changePassword() {
        this.password.open(this.user);
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
