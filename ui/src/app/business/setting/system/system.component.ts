import {Component, OnInit} from '@angular/core';
import {BaseModelComponent} from '../../../shared/class/BaseModelComponent';
import {System, SystemCreateRequest} from './system';
import {SystemService} from '../system.service';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import * as ipaddr from 'ipaddr.js';


@Component({
    selector: 'app-system',
    templateUrl: './system.component.html',
    styleUrls: ['./system.component.css']
})
export class SystemComponent extends BaseModelComponent<System> implements OnInit {

    items: System[] = [];
    item: SystemCreateRequest = new SystemCreateRequest();

    constructor(private systemService: SystemService, private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
        super(systemService);
    }

    ngOnInit(): void {
        this.listSystemSettings();
    }


    listSystemSettings() {
        this.systemService.list().subscribe(res => {
            this.items = res.items;
            this.item.vars = this.items[0].vars;
        });
    }

    onSubmit() {
        if (!this.checkIp()) {
            return;
        }
        this.systemService.create(this.item).subscribe(res => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    checkIp() {
        const ip = this.item.vars['ip'];
        if (!ipaddr.isValid(ip)) {
            this.commonAlertService.showAlert(this.translateService.instant('APP_IP_INVALID'), AlertLevels.ERROR);
            return false;
        }
        const ntp = this.item.vars['ntp_server'];
        if (ntp !== '' && !ipaddr.isValid(ntp)) {
            this.commonAlertService.showAlert(this.translateService.instant('APP_NTP_INVALID'), AlertLevels.ERROR);
            return false;
        }
        return true;
    }
}
