import {Component, OnInit} from '@angular/core';
import {MessageService} from './message.service';
import {System} from '../system/system';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../layout/common-alert/alert';

@Component({
    selector: 'app-message',
    templateUrl: './message.component.html',
    styleUrls: ['./message.component.css']
})
export class MessageComponent implements OnInit {

    item: System = new System();
    currentTab = 'EMAIL';
    loading = false;
    emailValid = false;

    constructor(private messageService: MessageService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
    }

    ngOnInit(): void {
        this.listByTab('EMAIL');
    }

    changeTab(tabName) {
        this.listByTab(tabName);
    }

    listByTab(tabName) {
        this.loading = false;
        this.messageService.getByTab(tabName).subscribe(res => {
            this.item = res;
            this.loading = false;
        }, error => {
            this.loading = false;
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    checkEmail() {
        this.emailValid = true;
    }

    onSubmit(item, tab) {
        this.messageService.postByTab(tab, item).subscribe(res => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    onCancel() {

    }
}
