import {Component, OnInit} from '@angular/core';
import {UserNotificationConfig} from '../message';
import {UserSubscribeService} from './user-subscribe.service';
import {SessionService} from '../../../shared/auth/session.service';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../layout/common-alert/alert';

@Component({
    selector: 'app-user-subscribe',
    templateUrl: './user-subscribe.component.html',
    styleUrls: ['./user-subscribe.component.css']
})
export class UserSubscribeComponent implements OnInit {

    loading = false;
    items: UserNotificationConfig[] = [];
    user;
    updateItem: UserNotificationConfig;

    constructor(private userSubscribeService: UserSubscribeService,
                private sessionService: SessionService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
    }

    ngOnInit(): void {
        this.loading = true;
        this.userSubscribeService.singleGet().subscribe(res => {
            this.loading = false;
            this.items = res;
        }, error => {
            this.loading = false;
        });
    }

    updateConfig(updateItem, type) {
        if (updateItem.vars[type] === 'DISABLE') {
            updateItem.vars[type] = 'ENABLE';
        } else {
            updateItem.vars[type] = 'DISABLE';
        }
        this.userSubscribeService.singleUpdate(updateItem).subscribe(res => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
