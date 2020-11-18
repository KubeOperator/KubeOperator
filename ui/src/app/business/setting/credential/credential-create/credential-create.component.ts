import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {CredentialCreateRequest} from '../credential';
import {AbstractControl, FormBuilder, FormGroup, NgForm, Validators} from '@angular/forms';
import {CredentialService} from '../credential.service';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {NamePattern, NamePatternHelper} from '../../../../constant/pattern';

@Component({
    selector: 'app-credential-create',
    templateUrl: './credential-create.component.html',
    styleUrls: ['./credential-create.component.css']
})
export class CredentialCreateComponent implements OnInit {

    namePattern = NamePattern;
    namePatternHelper = NamePatternHelper;
    opened = false;
    isSubmitGoing = false;
    item: CredentialCreateRequest = new CredentialCreateRequest();
    @ViewChild('credentialForm') credentialForm: NgForm;
    @Output() created = new EventEmitter();


    constructor(private service: CredentialService, private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private translateService: TranslateService) {
    }

    ngOnInit(): void {
    }

    open() {
        this.item = new CredentialCreateRequest();
        this.opened = true;
        this.item.type = 'password';
    }

    onCancel() {
        this.opened = false;
        this.item = new CredentialCreateRequest();
        this.item.type = 'password';
        this.credentialForm.resetForm(this.item);
    }

    onSubmit() {
        this.isSubmitGoing = true;
        this.service.create(this.item).subscribe(data => {
            this.opened = false;
            this.isSubmitGoing = false;
            this.created.emit();
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    validate() {
        return true;
    }

}


