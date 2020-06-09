import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {NgForm} from '@angular/forms';
import {LoginCredential} from './login-credential';
import {LoginService} from './login.service';
import {Router} from '@angular/router';
import {SessionService} from '../shared/session.service';
import {CommonRoutes} from '../globals';
import {TranslateService} from '@ngx-translate/core';

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


    constructor(private loginService: LoginService, private router: Router, private sessionService: SessionService,
                private translateService: TranslateService) {
    }

    ngOnInit(): void {
    }


    login() {
        this.loginService.login(this.loginCredential).subscribe(res => {
            this.isError = false;
            this.sessionService.cacheProfile(res);
            this.router.navigateByUrl(CommonRoutes.F2O_ROOT, {skipLocationChange: true}).then(r => console.log('login success'));
        }, error => this.handleError(error));
    }

    handleError(error: any) {
        this.isError = true;
        if (error.status === 504 || error.status === 502) {
            this.message = this.translateService.instant('APP_LOGIN_CONNECT_ERROR');
        } else if (error.status === 401) {
            this.message = this.translateService.instant('APP_LOGIN_CONNECT_CREDENTIAL_ERROR');
        } else {
            this.message = this.translateService.instant('APP_LOGIN_CONNECT_UNKNOWN_ERROR') + `${error.status}`;
        }
    }
}
