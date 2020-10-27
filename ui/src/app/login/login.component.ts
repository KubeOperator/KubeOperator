import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {NgForm} from '@angular/forms';
import {LoginCredential} from './login-credential';
import {LoginService} from './login.service';
import {Router} from '@angular/router';
import {SessionService} from '../shared/auth/session.service';
import {CommonRoutes} from '../constant/route';
import {TranslateService} from '@ngx-translate/core';
import {Theme} from '../business/setting/theme/theme';
import {ThemeService} from '../business/setting/theme/theme.service';
import {LicenseService} from '../business/setting/license/license.service';
import {VerificationCode} from '../shared/auth/session-user';

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
    verificationCode: VerificationCode = new VerificationCode();
    hasCode = false;

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
        this.checkLoginFailedNum();
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
        this.loginCredential.captchaId = this.verificationCode.captchaId;
        this.loginService.login(this.loginCredential).subscribe(res => {
            this.isError = false;
            this.sessionService.cacheProfile(res);
            localStorage.setItem('currentLanguage', this.loginCredential.language);
            this.translateService.use(this.loginCredential.language);
            this.router.navigateByUrl(CommonRoutes.KO_ROOT).then(r => console.log('login success'));
            localStorage.removeItem('loginErrorNum');
        }, error => this.handleError(error));
    }

    handleError(error: any) {
        this.isError = true;
        if (localStorage.getItem('loginErrorNum') != null) {
            const loginErrorNum = Number(localStorage.getItem('loginErrorNum'));
            if (loginErrorNum >= 3) {
                this.getVerificationCode();
            } else {
                const newNum = loginErrorNum + 1;
                localStorage.setItem('loginErrorNum', newNum.toString());
            }
        } else {
            localStorage.setItem('loginErrorNum', '1');
        }
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

    getVerificationCode() {
        this.loginService.getCode().subscribe(res => {
            this.verificationCode = res;
            this.hasCode = true;
        }, error => {
            this.message = this.translateService.instant(error.msg);
        });
    }

    checkLoginFailedNum() {
        if (localStorage.getItem('loginErrorNum') != null) {
            const loginErrorNum = Number(localStorage.getItem('loginErrorNum'));
            if (loginErrorNum >= 3) {
                this.getVerificationCode();
            }
        }
    }
}
