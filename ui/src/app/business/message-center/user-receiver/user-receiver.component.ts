import {Component, OnInit} from '@angular/core';
import {UserReceiverService} from './user-receiver.service';
import {UserReceiver} from '../message';
import {SessionService} from '../../../shared/auth/session.service';
import {SystemService} from '../../setting/system.service';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../layout/common-alert/alert';

@Component({
    selector: 'app-user-receiver',
    templateUrl: './user-receiver.component.html',
    styleUrls: ['./user-receiver.component.css']
})
export class UserReceiverComponent implements OnInit {

    item: UserReceiver = new UserReceiver();
    submitGoing = false;
    user;

    constructor(private userReceiverService: UserReceiverService,
                private sessionService: SessionService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
    }

    ngOnInit(): void {
        this.userReceiverService.singleGet().subscribe(res => {
            this.item = res;
        }, error => {
        });
    }

    onCancel() {

    }

    onSubmit() {
        this.userReceiverService.singleUpdate(this.item).subscribe(res => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
