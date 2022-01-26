import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {NgForm} from '@angular/forms';
import {LoginCredential} from './login-credential';
import {Router} from '@angular/router';
import {SessionService} from '../shared/auth/session.service';
import {CommonRoutes} from '../constant/route';
import {TranslateService} from '@ngx-translate/core';
import {Theme} from '../business/setting/theme/theme';
import {ThemeService} from '../business/setting/theme/theme.service';
import {Captcha} from '../shared/auth/session-user';
import {ResetPasswordComponent} from './reset-password/reset-password.component';

@Component({
    selector: 'app-login',
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

    @ViewChild('loginForm', {static: true}) loginForm: NgForm;
    @Input() loginCredential: LoginCredential = new LoginCredential();
    @ViewChild(ResetPasswordComponent, {static: true}) resetPwdDialog: ResetPasswordComponent;
    message: string;
    isError = false;
    isFirst = false;
    theme: Theme;
    captcha: Captcha = new Captcha();

    constructor(private router: Router,
                private themeService: ThemeService,
                private sessionService: SessionService,
                private translateService: TranslateService) {
    }

    ngOnInit(): void {
        const currentLanguage = localStorage.getItem('currentLanguage');
        if (currentLanguage) {
            this.loginCredential.language = currentLanguage;
        } else {
            this.loginCredential.language = 'zh-CN';
        }
        this.loadTheme();
        this.createCaptcha();
    }

    reset() {
        this.loginForm.resetForm();
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
        this.loginCredential.captchaId = this.captcha.captchaId;
        this.sessionService.login(this.loginCredential).subscribe(res => {
            this.isError = false;
            if (res.user.isFirst) {
                this.resetOpen()
            } else {
              this.router.navigateByUrl(CommonRoutes.KO_ROOT);
            }
            this.sessionService.cacheProfile(res);
            localStorage.setItem('currentLanguage', this.loginCredential.language);
            this.translateService.use(this.loginCredential.language);
        }, error => this.handleError(error));
    }

    resetOpen() {
        this.resetPwdDialog.open(this.loginCredential.username, this.loginCredential.password);
    }

    handleError(error: any) {
        this.isError = true
        switch (error.status) {
            case 504:
                this.message = this.translateService.instant('APP_LOGIN_CONNECT_ERROR');
                break;
            case 400:
                this.message = error.error.msg;
                break;
            default:
                this.message = this.translateService.instant('APP_LOGIN_CONNECT_UNKNOWN_ERROR') + `${error.status}`;
        }
        this.createCaptcha();
    }

    createCaptcha() {
        this.sessionService.getCode().subscribe(res => {
            this.captcha = res;
        }, error => {
            this.message = this.translateService.instant(error.msg);
        });
    }
}
