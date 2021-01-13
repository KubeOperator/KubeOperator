import {Component, OnInit} from '@angular/core';
import {System, SystemCreateRequest} from '../system/system';
import {SystemService} from '../system.service';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../layout/common-alert/alert';

@Component({
    selector: 'app-email',
    templateUrl: './email.component.html',
    styleUrls: ['./email.component.css']
})
export class EmailComponent implements OnInit {


    item: System = new System();
    valid = false;
    createItem: SystemCreateRequest = new SystemCreateRequest();
    checking = false;

    constructor(private systemService: SystemService, private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
    }

    ngOnInit(): void {
        this.getSetting();
    }


    checkValid(type) {
        this.createItem.vars = this.item.vars;
        this.createItem.tab = 'EMAIL';
        this.checking = true;
        this.systemService.checkBy(type, this.createItem).subscribe(res => {
            this.valid = true;
            this.checking = false;
            this.commonAlertService.showAlert(this.translateService.instant('APP_CHECK_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.valid = false;
            this.checking = false;
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    onSubmit() {
        this.createItem.vars = this.item.vars;
        this.createItem.tab = 'EMAIL';
        this.systemService.create(this.createItem).subscribe(res => {
            let key = 'APP_ADD_SUCCESS';
            if (this.item !== undefined && this.item.tab !== '') {
                key = 'APP_UPDATE_SUCCESS';
            }
            this.commonAlertService.showAlert(this.translateService.instant(key), AlertLevels.SUCCESS);
            this.getSetting();
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    changeStatus(status) {
        this.valid = status !== 'ENABLE';
    }

    getSetting() {
        this.systemService.getByTab('EMAIL').subscribe(res => {
            this.item = res;
        });
    }
}
