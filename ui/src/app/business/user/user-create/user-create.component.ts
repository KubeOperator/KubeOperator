import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {User, UserCreateRequest} from '../user';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {UserService} from '../user.service';
import {NgForm} from '@angular/forms';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../shared/common-component/modal-alert/modal-alert.service';
import {NamePattern, PasswordPattern} from '../../../constant/pattern';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';


@Component({
    selector: 'app-user-create',
    templateUrl: './user-create.component.html',
    styleUrls: ['./user-create.component.css']
})
export class UserCreateComponent extends BaseModelDirective<User> implements OnInit {

    opened = false;
    isSubmitGoing = false;
    item: UserCreateRequest = new UserCreateRequest();
    passwordPattern = PasswordPattern;
    namePattern = NamePattern;
    @ViewChild('userForm') userForm: NgForm;
    @Output() created = new EventEmitter();

    private validationStateMap: any = {
        password: true,
        rePassword: true,
        namePwd: true,
        rePwdCheck: true,
    };

    constructor(private userService: UserService, private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private translateService: TranslateService) {
        super(userService);
    }

    ngOnInit(): void {
    }

    public get isValid(): boolean {
        if (this.userForm && this.userForm.form.get('password')) {
            return this.userForm.valid && (this.userForm.form.get('password').value === this.userForm.form.get('rePassword').value);
        }
        return false;
    }

    open() {
        this.opened = true;
        this.item = new UserCreateRequest();
        this.userForm.resetForm();
    }

    onCancel() {
        this.opened = false;
    }

    getValidationState(key: string): boolean {
        return this.validationStateMap[key];
    }

    handleValidation(key) {
        const cont = this.userForm.controls[key];
        if (cont && cont.invalid && !cont.hasError) {
            this.validationStateMap[key] = false;
            return;
        }
        this.validationStateMap[key] = true;
        if (this.userForm.form.get('password').value === this.userForm.form.get('name').value) {
            this.userForm.controls['password'].setErrors({namePwdError: false});
            this.validationStateMap['namePwd'] = false;
            return;
        } else {
            this.validationStateMap['namePwd'] = true;
        }
        const r = /^(?=.*\d)(?=.*[a-zA-Z])[\da-zA-Z~!@#$%^&*]{6,30}$/g;
        r.lastIndex = 0;
        if (!r.test(this.userForm.form.get('password').value)) {
            this.userForm.controls['password'].setErrors({passwordError: false});
            this.validationStateMap['password'] = false;
            return;
        } else {
            this.validationStateMap['password'] = true;
        }

        if (this.userForm.form.get('rePassword').value !== null && this.userForm.form.get('password').value !== this.userForm.form.get('rePassword').value) {
            this.userForm.controls[key].setErrors({rePwdError: false});
            this.validationStateMap['rePwdCheck'] = false;
        } else {
            this.validationStateMap['rePwdCheck'] = true;
            this.userForm.controls['password'].setErrors(null);
            this.userForm.controls['rePassword'].setErrors(null);
        }
    }

    onSubmit() {
        if (this.userForm.invalid || !this.isValid) {
            return;
        }
        if (this.item.name === this.item.password) {
            this.modalAlertService.showAlert(this.translateService.instant('USERNAME_PWD_INVALID'), AlertLevels.ERROR);
            return;
        }
        this.isSubmitGoing = true;
        this.userService.create(this.item).subscribe(data => {
            this.opened = false;
            this.isSubmitGoing = false;
            this.created.emit();
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
