import {Component, OnInit} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {Registry} from '../registry';
import {SystemService} from '../../system.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {RegistryService} from '../registry.service';
import {System} from '../../system/system';

@Component({
    selector: 'app-registry-list',
    templateUrl: './registry-list.component.html',
    styleUrls: ['./registry-list.component.css']
})
export class RegistryListComponent extends BaseModelDirective<Registry> implements OnInit {
    systemItem: System = new System();

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
        if (this.systemItem.vars['arch_type'] === 'Mixed') {
            console.log('mixed running');
            // this.registryService.mixedGet().subscribe(res => {
            //     this.items = res;
            // });
        }
    }
}
