import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {User} from '../user';
import {UserService} from '../user.service';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {SessionService} from '../../../shared/auth/session.service';
import {SessionUser} from '../../../shared/auth/session-user';

@Component({
    selector: 'app-user-delete',
    templateUrl: './user-delete.component.html',
    styleUrls: ['./user-delete.component.css']
})
export class UserDeleteComponent extends BaseModelDirective<User> implements OnInit {

    opened = false;
    items: User[] = [];
    user: SessionUser = new SessionUser();

    @Output()
    deleted = new EventEmitter();

    constructor(private userService: UserService, private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private translateService: TranslateService,
                private sessionService: SessionService) {
        super(userService);
    }

    ngOnInit(): void {
        this.sessionService.getProfile().subscribe(res => {
            this.user = res.user;
        })
    }

    open(items) {
        this.items = items;
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        let check = true;
        for (const item of this.items) {
            if (item.name === this.user.name) {
                check = false;
                break;
            }
        }
        if (!check) {
            this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_US'), AlertLevels.ERROR);
            this.opened = false;
            return;
        }
        this.service.batch('delete', this.items).subscribe(data => {
            this.deleted.emit();
            this.opened = false;
            this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
            this.opened = false;
        });
    }
}
