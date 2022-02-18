import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {CredentialService} from '../credential.service';
import {Credential} from '../credential';
import {NgForm} from '@angular/forms';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-credential-edit',
    templateUrl: './credential-edit.component.html',
    styleUrls: ['./credential-edit.component.css']
})
export class CredentialEditComponent implements OnInit {

    item: Credential = new Credential();
    opened = false;
    isSubmitGoing = false;
    @ViewChild('credentialEditForm') credentialForm: NgForm;
    @Output() edit = new EventEmitter();

    constructor(private service: CredentialService, private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private translateService: TranslateService) {
    }

    ngOnInit(): void {
    }

    open(item: Credential) {
        Object.assign(this.item, item);
        this.item.password = '';
        this.opened = true;
    }

    onCancel() {
        this.item = new Credential();
        this.credentialForm.resetForm(this.item);
        this.opened = false;
    }

    onSubmit() {
        this.isSubmitGoing = true;
        this.service.update(this.item.name, this.item).subscribe(data => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
            this.opened = false;
            this.isSubmitGoing = false;
            this.edit.emit();
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.msg ? error.msg : error.error.msg, AlertLevels.ERROR);
        });
    }
}
