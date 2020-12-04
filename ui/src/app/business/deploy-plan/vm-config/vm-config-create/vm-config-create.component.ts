import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {VmConfig} from '../vm-config';
import {VmConfigService} from '../vm-config.service';
import {NgForm} from '@angular/forms';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {
    VmConfigPattern
} from '../../../../constant/pattern';

@Component({
    selector: 'app-vm-config-create',
    templateUrl: './vm-config-create.component.html',
    styleUrls: ['./vm-config-create.component.css']
})
export class VmConfigCreateComponent extends BaseModelDirective<VmConfig> implements OnInit {


    opened = false;
    item: VmConfig = new VmConfig();
    isSubmitGoing = false;
    @ViewChild('vmConfigForm', {static: true}) vmConfigForm: NgForm;
    @Output() created = new EventEmitter();
    namePattern = VmConfigPattern;

    constructor(private vmConfigService: VmConfigService,
                private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
        super(vmConfigService);
    }

    ngOnInit(): void {
    }

    open() {
        this.opened = true;
        this.item = new VmConfig();
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.vmConfigService.create(this.item).subscribe(res => {
            this.opened = false;
            this.created.emit();
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
