import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {VmConfig} from '../vm-config';
import {VmConfigService} from '../vm-config.service';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {NgForm} from '@angular/forms';
import {AlertLevels} from '../../../../layout/common-alert/alert';

@Component({
    selector: 'app-vm-config-update',
    templateUrl: './vm-config-update.component.html',
    styleUrls: ['./vm-config-update.component.css']
})
export class VmConfigUpdateComponent extends BaseModelDirective<VmConfig> implements OnInit {

    opened = false;
    item: VmConfig = new VmConfig();
    @Output() updated = new EventEmitter();
    @ViewChild('vmConfigEditForm', {static: true}) vmConfigEditForm: NgForm;

    constructor(private vmConfigService: VmConfigService,
                private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
        super(vmConfigService);
    }

    ngOnInit(): void {
    }

    open(item) {
        this.opened = true;
        Object.assign(this.item, item);
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.vmConfigService.update(this.item.name, this.item).subscribe(data => {
            this.onCancel();
            this.updated.emit();
            this.commonAlertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
