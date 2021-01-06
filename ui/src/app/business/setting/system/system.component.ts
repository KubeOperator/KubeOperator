import {Component, OnInit} from '@angular/core';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
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
export class SystemComponent extends BaseModelDirective<System> implements OnInit {

    item: System = new System();
    createItem: SystemCreateRequest = new SystemCreateRequest();

    constructor(private systemService: SystemService, private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
        super(systemService);
    }

    ngOnInit(): void {
        this.listSystemSettings();
    }

    listSystemSettings() {
        this.systemService.singleGet().subscribe(res => {
            this.item = res;
            if ( this.item.vars['REGISTRY_PROTOCOL'] === undefined || this.item.vars['REGISTRY_PROTOCOL'] === '' ) {
                    this.item.vars['REGISTRY_PROTOCOL'] = 'http';
                }
        });
    }

    onSubmit() {
        if (!this.checkIp()) {
            return;
        }
        if (this.item.vars['ip'] === this.item.vars['ntp_server']) {
            this.commonAlertService.showAlert(this.translateService.instant('APP_SYSTEM_EQUAL'), AlertLevels.ERROR);
            return;
        }

        this.createItem.vars = this.item.vars;
        this.systemService.create(this.createItem).subscribe(res => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
            window.location.reload();
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
        return true;
    }
}
