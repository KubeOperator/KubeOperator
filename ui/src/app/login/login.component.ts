import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {NgForm} from '@angular/forms';
import {LoginCredential} from './login-credential';
import {LoginService} from './login.service';
import {Router} from '@angular/router';
import {SessionService} from '../shared/auth/session.service';
import {CommonRoutes} from '../constant/route';
import {TranslateService} from '@ngx-translate/core';
import {Theme} from "../business/setting/theme/theme";
import {ThemeService} from "../business/setting/theme/theme.service";
import {LicenseService} from "../business/setting/license/license.service";

@Component({
    selector: 'app-login',
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

    @ViewChild('loginForm', {static: true}) loginForm: NgForm;
    @Input() loginCredential: LoginCredential = new LoginCredential();
    message: string;
    isError = false;
    theme: Theme;

    constructor(private loginService: LoginService,
                private router: Router,
                private themeService: ThemeService,
                private sessionService: SessionService,
                private translateService: TranslateService,
                private licenseService: LicenseService) {
    }

    ngOnInit(): void {
        const currentLanguage = localStorage.getItem('currentLanguage');
        if (currentLanguage) {
            this.loginCredential.language = currentLanguage;
        } else {
            this.loginCredential.language = 'zh-CN';
        }
        this.loadTheme();
    }

    loadTheme() {
        this.themeService.get().subscribe(data => {
            this.theme = data;
            if (this.theme.systemName) {
                document.title = this.theme.systemName;
            }
        });
    }
    login() {
        this.loginService.login(this.loginCredential).subscribe(res => {
            this.isError = false;
            this.sessionService.cacheProfile(res);
            localStorage.setItem('currentLanguage', this.loginCredential.language);
            this.translateService.use(this.loginCredential.language);
            this.router.navigateByUrl(CommonRoutes.KO_ROOT).then(r => console.log('login success'));
        }, error => this.handleError(error));
    }

    handleError(error: any) {
        this.isError = true;
        switch (error.status) {
            case 500:
                this.message = error.error.msg;
                break;
            case 504:
                this.message = this.translateService.instant('APP_LOGIN_CONNECT_ERROR');
                break;
            case 401:
                this.message = error.error.msg;
                break;
            default:
                this.message = this.translateService.instant('APP_LOGIN_CONNECT_UNKNOWN_ERROR') + `${error.status}`;
        }
    }
}
