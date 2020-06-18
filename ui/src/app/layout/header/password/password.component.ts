import {Component, OnInit, ViewChild} from '@angular/core';
import {PasswordPattern} from '../../../constant/pattern';
import {NgForm} from '@angular/forms';
import {UserService} from '../../../business/user/user.service';
import {SessionUser} from '../../../shared/auth/session-user';
import {ChangePasswordRequest} from '../../../business/user/user';
import {ModalAlertService} from '../../../shared/common-component/modal-alert/modal-alert.service';
import {AlertLevels} from '../../common-alert/alert';

@Component({
    selector: 'app-password',
    templateUrl: './password.component.html',
    styleUrls: ['./password.component.css']
})
export class PasswordComponent implements OnInit {

    opened = false;
    password: string;
    confirmPassword: string;
    original: string;
    submitGoing = false;
    passwordPattern = PasswordPattern;
    user: SessionUser = new SessionUser();
    changePasswordRequest: ChangePasswordRequest = new ChangePasswordRequest();
    @ViewChild('passForm', {static: true}) passForm: NgForm;

    constructor(private userService: UserService, private modalAlertService: ModalAlertService) {
    }

    ngOnInit(): void {
    }

    open(user) {
        this.user = user;
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
        this.passForm.resetForm();
    }

    onSubmit() {
        this.submitGoing = true;
        this.changePasswordRequest = {
            password: this.password,
            original: this.original,
            name: this.user.name
        };

        this.userService.changePassword(this.changePasswordRequest).subscribe(res => {
            this.submitGoing = false;
            this.opened = false;
        }, error => {
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
            this.submitGoing = false;
        });
    }

    checkPassword() {
        return this.password === this.confirmPassword;
    }
}
