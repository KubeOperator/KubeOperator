import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {User} from '../user';
import {UserService} from '../user.service';
import {NgForm} from '@angular/forms';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-user-update',
    templateUrl: './user-update.component.html',
    styleUrls: ['./user-update.component.css']
})
export class UserUpdateComponent extends BaseModelDirective<User> implements OnInit {

    opened = false;
    item: User = new User();
    itemEmail: string = '';
    isSubmitGoing = false;
    @ViewChild('userForm') userFrom: NgForm;

    @Output()
    update = new EventEmitter();

    constructor(private userService: UserService, private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private translateService: TranslateService) {
        super(userService);
    }

    ngOnInit(): void {
    }

    open(item) {
        if (item.type === 'LDAP') {
            this.commonAlertService.showAlert(this.translateService.instant('APP_USER_LDAP_UPDATE_ERROR'), AlertLevels.ERROR);
            return;
        }
        this.opened = true;
        Object.assign(this.item, item);
    }

    onCancel() {
        this.itemEmail = '';
        this.opened = false;
    }

    onSubmit() {
        this.item.email = this.itemEmail;
        this.isSubmitGoing = true;
        this.userService.update(this.item.name, this.item).subscribe(data => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
            this.opened = false;
            this.isSubmitGoing = false;
            this.update.emit();
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
        this.itemEmail = '';
    }
}
