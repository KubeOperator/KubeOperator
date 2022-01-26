import {Component, OnInit, ViewChild} from '@angular/core';
import {NgForm} from '@angular/forms';
import {Router} from '@angular/router';
import {UserService} from '../../business/user/user.service';
import {ModalAlertService} from '../../shared/common-component/modal-alert/modal-alert.service';
import {AlertLevels} from '../../layout/common-alert/alert';
import {CommonRoutes} from '../../constant/route';
import {PasswordPattern} from '../../constant/pattern';
import {ChangePasswordRequest} from '../../business/user/user';

@Component({
    selector: 'app-reset-password',
    templateUrl: './reset-password.component.html',
    styleUrls: ['./reset-password.component.css']
})
export class ResetPasswordComponent implements OnInit {

    opened = false;
    
    loading = false;
    password: string;
    confirmPassword: string;
    passwordPattern = PasswordPattern;
    email: string;
    changePasswordRequest: ChangePasswordRequest = new ChangePasswordRequest();

    @ViewChild('resetPasswordFrom', {static: true}) resetPwdForm: NgForm;

    constructor(private router: Router,
                private userService: UserService,
                private modalAlertService: ModalAlertService) {
    }

    ngOnInit(): void {
    }

    open(name, original) {
        this.opened = true;
        this.changePasswordRequest.name = name;
        this.changePasswordRequest.original = original;
        this.resetPwdForm.resetForm();
    }

    close() {
        this.opened = false;
    }

    reset() {
        this.changePasswordRequest.password = this.password
        this.userService.changePassword(this.changePasswordRequest).subscribe(res => {
            this.router.navigateByUrl(CommonRoutes.KO_ROOT);
            this.opened = false;
        }, error => {
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    checkPassword() {
        return this.password === this.confirmPassword;
    }
}
