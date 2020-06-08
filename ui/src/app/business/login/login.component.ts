import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {NgForm} from '@angular/forms';
import {LoginCredential} from './login-credential';
import {LoginService} from './login.service';
import {Router} from '@angular/router';

@Component({
    selector: 'app-login',
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

    @ViewChild('loginForm', {static: true}) loginForm: NgForm;
    @Input() loginCredential: LoginCredential = new LoginCredential();
    message: string;

    constructor(private loginService: LoginService, private router: Router) {
    }

    ngOnInit(): void {
    }


    login() {
        this.loginService.login(this.loginCredential).subscribe(res => {
            this.router.navigateByUrl('/', {skipLocationChange: true}).then(r => console.log('login success'));
        }, error => this.handleError(error));
    }

    handleError(error: any) {
        // this.signInStatus = signInStatusError;
        if (error.status === 504 || error.status === 502) {
            this.message = 'kubeOperator Api 连接失败！';
        } else if (error.status === 400) {
            this.message = '用户名或密码错误！';
        } else {
            this.message = `未知错误,code: ${error.status}`;
        }
    }

}
