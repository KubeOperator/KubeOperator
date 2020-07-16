import {Component, OnInit} from '@angular/core';
import {BaseModelComponent} from '../../../shared/class/BaseModelComponent';
import {System, SystemCreateRequest} from './system';
import {SystemService} from '../system.service';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';

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
        this.systemService.create(this.item).subscribe(res => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.SUCCESS);
        });
    }
}
