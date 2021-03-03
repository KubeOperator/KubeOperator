import {Component, OnInit} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {Registry} from '../registry';
import {SystemService} from '../../system.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {RegistryService} from '../registry.service';
import {System, SystemCreateRequest} from '../../system/system';
import {AlertLevels} from '../../../../layout/common-alert/alert';

@Component({
    selector: 'app-registry-list',
    templateUrl: './registry-list.component.html',
    styleUrls: ['./registry-list.component.css']
})
export class RegistryListComponent extends BaseModelDirective<Registry> implements OnInit {
    item: System = new System();
    systemItem: System = new System();
    createItem: SystemCreateRequest = new SystemCreateRequest();

    constructor(private systemService: SystemService, private commonAlertService: CommonAlertService,
                private translateService: TranslateService, private registryService: RegistryService) {
        super(registryService);
    }

    ngOnInit(): void {
        super.ngOnInit();
        this.listSystemSettings();
    }

    listSystemSettings() {
        this.systemService.singleGet().subscribe(res => {
            this.systemItem = res;
        });
    }

    SingleOnSubmit() {
        this.createItem.vars = this.systemItem.vars;
        this.createItem.tab = 'SYSTEM';
        this.systemService.create(this.createItem).subscribe(res => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
