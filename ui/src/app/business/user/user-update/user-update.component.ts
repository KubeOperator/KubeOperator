import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelComponent} from '../../../shared/class/BaseModelComponent';
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
export class UserUpdateComponent extends BaseModelComponent<User> implements OnInit {

    opened = false;
    item: User = new User();
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
        this.opened = true;
        Object.assign(this.item, item);
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.isSubmitGoing = true;
        this.userService.update(this.item.name, this.item).subscribe(data => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
            this.opened = false;
            this.isSubmitGoing = false;
            this.update.emit();
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.msg, AlertLevels.ERROR);
        });
    }
}
