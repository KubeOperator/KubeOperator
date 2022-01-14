import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Registry, RegistryCreateRequest} from '../registry';
import {NgForm} from '@angular/forms';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {RegistryService} from '../registry.service';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {IpPattern} from '../../../../constant/pattern';

@Component({
    selector: 'app-registry-create',
    templateUrl: './registry-create.component.html',
    styleUrls: ['./registry-create.component.css']
})
export class RegistryCreateComponent extends BaseModelDirective<Registry> implements OnInit {
    opened = false;
    isSubmitGoing = false;
    ipPattern = IpPattern;

    item: RegistryCreateRequest = new RegistryCreateRequest();
    registries: Registry[] = [];
    @ViewChild('registryForm') registryForm: NgForm;
    @Output() created = new EventEmitter();

    constructor(private registryService: RegistryService, private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private translateService: TranslateService) {
        super(registryService);
    }

    ngOnInit(): void {
    }

    open() {
        this.registryService.list().subscribe(data => {
            this.registries = data.items;
        });
        this.opened = true;
        this.item = new RegistryCreateRequest();
        this.setDefaultValue();
        // this.registryForm.resetForm();
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.isSubmitGoing = true;
        this.registryService.create(this.item).subscribe(data => {
            this.opened = false;
            this.isSubmitGoing = false;
            this.created.emit();
            this.onCancel();
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
            window.location.reload();
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    setDefaultValue() {
        this.item.architecture = 'x86_64';
        this.item.protocol = 'http';
    }
}
