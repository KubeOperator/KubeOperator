import {Component, OnInit} from '@angular/core';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {User} from '../user';
import {UserService} from '../user.service';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-user-list',
    templateUrl: './user-list.component.html',
    styleUrls: ['./user-list.component.css']
})
export class UserListComponent extends BaseModelDirective<User> implements OnInit {

    constructor(private userService: UserService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
        super(userService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }

    updateUser(item) {
        this.userService.update(item.name, item).subscribe(data => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
            this.refresh();
        });
    }
}
