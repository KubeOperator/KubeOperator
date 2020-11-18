import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {User, UserCreateRequest} from '../user';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {UserService} from '../user.service';
import {NgForm} from '@angular/forms';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../shared/common-component/modal-alert/modal-alert.service';
import {NamePattern, NamePatternHelper, PasswordPattern} from '../../../constant/pattern';
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
    namePatternHelper = NamePatternHelper;
    isPasswordMatch = false;
    @ViewChild('userForm') userForm: NgForm;
    @Output() created = new EventEmitter();

    constructor(private userService: UserService, private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private translateService: TranslateService) {
        super(userService);
    }

    ngOnInit(): void {
    }

    open() {
        this.opened = true;
        this.item = new UserCreateRequest();
        this.userForm.resetForm();
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
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

    checkPassword() {
        this.isPasswordMatch = this.item.password === this.item.confirmPassword;
    }
}
