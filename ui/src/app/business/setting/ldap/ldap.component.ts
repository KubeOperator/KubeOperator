import {Component, OnInit} from '@angular/core';
import {SystemService} from '../system.service';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {System, SystemCreateRequest} from '../system/system';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {LdapService} from './ldap.service';

@Component({
    selector: 'app-ldap',
    templateUrl: './ldap.component.html',
    styleUrls: ['./ldap.component.css']
})
export class LdapComponent extends BaseModelDirective<System> implements OnInit {

    item: System = new System();
    createItem: SystemCreateRequest = new SystemCreateRequest();

    constructor(private systemService: SystemService, private commonAlertService: CommonAlertService,
                private translateService: TranslateService, private ldapService: LdapService) {
        super(systemService);
    }

    ngOnInit(): void {
        this.listSystemSettings();
    }


    listSystemSettings() {
        this.systemService.singleGet().subscribe(res => {
            this.item = res;
        });
    }

    onSubmit() {
        this.createItem.vars = this.item.vars;
        this.createItem.tab = 'LDAP';
        this.ldapService.ldapCreate(this.createItem).subscribe(res => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    onSync() {
        this.createItem.vars = this.item.vars;
        this.createItem.tab = 'LDAP';
        this.ldapService.ldapSync(this.createItem).subscribe(res => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_SYNC_NOTE'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
