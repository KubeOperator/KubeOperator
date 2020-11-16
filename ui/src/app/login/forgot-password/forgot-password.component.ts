import {Component, OnInit, ViewChild} from '@angular/core';
import {NgForm} from '@angular/forms';
import {ResetPassword} from './reset-password';
import {UserService} from '../../business/user/user.service';
import {ModalAlertService} from '../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../layout/common-alert/alert';

@Component({
    selector: 'app-forgot-password',
    templateUrl: './forgot-password.component.html',
    styleUrls: ['./forgot-password.component.css']
})
export class ForgotPasswordComponent implements OnInit {

    opened = false;
    loading = false;
    email: string;
    item: ResetPassword = new ResetPassword();
    @ViewChild('forgotPasswordFrom', {static: true}) forgotPwdForm: NgForm;

    constructor(private userService: UserService,
                private modalAlertService: ModalAlertService,
                private translateService: TranslateService) {
    }

    ngOnInit(): void {
    }

    open() {
        this.opened = true;
        this.item = new ResetPassword();
        this.forgotPwdForm.resetForm();
    }

    close() {
        this.opened = false;
    }

    send() {
        this.userService.resetPassword(this.item).subscribe(res => {
            this.modalAlertService.showAlert(this.translateService.instant('RESET_PWD_SUCCESS'), AlertLevels.SUCCESS);
            this.opened = false;
        }, error => {
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
