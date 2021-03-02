import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Credential} from '../../credential/credential';
import {Registry} from '../registry';
import {NgForm} from '@angular/forms';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {RegistryService} from '../registry.service';

@Component({
    selector: 'app-registry-update',
    templateUrl: './registry-update.component.html',
    styleUrls: ['./registry-update.component.css']
})
export class RegistryUpdateComponent implements OnInit {

    item = new Registry();
    opened = false;
    isSubmitGoing = false;

    @ViewChild('registryUpdateForm') registryForm: NgForm;
    @Output() update = new EventEmitter();

    constructor(private service: RegistryService, private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private translateService: TranslateService) {
    }

    ngOnInit(): void {
    }

    open(item: Credential) {
        Object.assign(this.item, item);
        this.opened = true;
    }

    onCancel() {
        this.item = new Registry();
        this.registryForm.resetForm(this.item);
        this.opened = false;
    }

    onSubmit() {
        this.isSubmitGoing = true;
        this.service.updateRegistryBy(this.item.architecture, this.item).subscribe(data => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
            this.opened = false;
            this.isSubmitGoing = false;
            this.update.emit();
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.msg, AlertLevels.ERROR);
        });
    }
}
