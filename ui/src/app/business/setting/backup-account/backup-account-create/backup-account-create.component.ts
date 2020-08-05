import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {BackupAccount, BackupAccountCreateRequest} from '../backup-account';
import {BackupAccountService} from '../backup-account.service';
import {NgForm} from '@angular/forms';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../layout/common-alert/alert';


@Component({
    selector: 'app-backup-account-create',
    templateUrl: './backup-account-create.component.html',
    styleUrls: ['./backup-account-create.component.css']
})
export class BackupAccountCreateComponent extends BaseModelComponent<BackupAccount> implements OnInit {

    opened = false;
    isSubmitGoing = false;
    item: BackupAccountCreateRequest = new BackupAccountCreateRequest();
    buckets = [];
    @Output() created = new EventEmitter();
    @ViewChild('backupAccountForm') backupAccountForm: NgForm;


    constructor(private backupAccountService: BackupAccountService, private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private translateService: TranslateService) {
        super(backupAccountService);
    }

    ngOnInit(): void {
    }

    open() {
        this.opened = true;
        this.item = new BackupAccountCreateRequest();
        this.buckets = [];
    }

    changeType() {

    }

    getBuckets() {
        this.backupAccountService.listBuckets(this.item).subscribe(res => {
            this.buckets = res;
        }, error => {

        });
    }

    onCancel() {
        this.opened = false;
        this.item = new BackupAccountCreateRequest();
        this.backupAccountForm.resetForm(this.item);
    }

    onSubmit() {
        this.isSubmitGoing = true;
        this.backupAccountService.create(this.item).subscribe(res => {
            this.isSubmitGoing = false;
            this.created.emit();
            this.onCancel();
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
