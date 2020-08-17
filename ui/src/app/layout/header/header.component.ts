import {Component, OnInit, ViewChild} from '@angular/core';
import {SessionService} from '../../shared/auth/session.service';
import {SessionUser} from '../../shared/auth/session-user';
import {Router} from '@angular/router';
import {CommonRoutes} from '../../constant/route';
import {PasswordComponent} from './password/password.component';
import {AboutComponent} from "./about/about.component";

@Component({
    selector: 'app-header',
    templateUrl: './header.component.html',
    styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit {

    user: SessionUser = new SessionUser();

    @ViewChild(PasswordComponent, {static: true})
    password: PasswordComponent;

    @ViewChild(AboutComponent, {static: true})
    about: AboutComponent;

    constructor(private sessionService: SessionService, private router: Router) {
    }

    ngOnInit(): void {
        this.getProfile();
    }

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

    logOut() {
        this.sessionService.clear();
        this.router.navigateByUrl(CommonRoutes.LOGIN).then();
    }

    openAbout() {
        this.about.open();
    }
}
