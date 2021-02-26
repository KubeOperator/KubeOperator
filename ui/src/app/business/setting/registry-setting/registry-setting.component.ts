import {Component, OnInit} from '@angular/core';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {System, SystemCreateRequest} from '../system/system';
import {SystemService} from '../system.service';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {Registry, RegistryCreateRequest} from './registry';

@Component({
    selector: 'app-registry-setting',
    templateUrl: './registry-setting.component.html',
    styleUrls: ['./registry-setting.component.css']
})
export class RegistrySettingComponent extends BaseModelDirective<System> implements OnInit {

    item: System = new System();
    createItem: SystemCreateRequest = new SystemCreateRequest();
    mixedItem: Registry = new Registry();
    mixedCreateItem: RegistryCreateRequest = new RegistryCreateRequest();
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
        });
    }

    SingleOnSubmit() {
        this.createItem.vars = this.item.vars;
        this.createItem.tab = 'SYSTEM';
        this.systemService.create(this.createItem).subscribe(res => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
            window.location.reload();
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    MixedOnSubmit() {
        console.log('Mixed');
    }
}
