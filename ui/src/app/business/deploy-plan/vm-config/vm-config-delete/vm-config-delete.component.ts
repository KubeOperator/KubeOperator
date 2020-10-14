import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {VmConfig} from '../vm-config';
import {VmConfigService} from '../vm-config.service';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../layout/common-alert/alert';

@Component({
    selector: 'app-vm-config-delete',
    templateUrl: './vm-config-delete.component.html',
    styleUrls: ['./vm-config-delete.component.css']
})
export class VmConfigDeleteComponent extends BaseModelDirective<VmConfig> implements OnInit {

    opened = false;
    items;
    @Output() deleted = new EventEmitter();

    constructor(private vmConfigService: VmConfigService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
        super(vmConfigService);
    }

    ngOnInit(): void {
    }

    open(items) {
        this.opened = true;
        this.items = items;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.vmConfigService.batch('delete', this.items).subscribe(data => {
            this.deleted.emit();
            this.opened = false;
            this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
            this.opened = false;
        });
    }
}
