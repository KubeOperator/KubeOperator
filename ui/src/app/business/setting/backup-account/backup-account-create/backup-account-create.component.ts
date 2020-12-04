import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {BackupAccount, BackupAccountCreateRequest} from '../backup-account';
import {BackupAccountService} from '../backup-account.service';
import {NgForm} from '@angular/forms';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {NamePattern} from '../../../../constant/pattern';


@Component({
    selector: 'app-backup-account-create',
    templateUrl: './backup-account-create.component.html',
    styleUrls: ['./backup-account-create.component.css']
})
export class BackupAccountCreateComponent extends BaseModelDirective<BackupAccount> implements OnInit {

    namePattern = NamePattern;
    opened = false;
    isSubmitGoing = false;
    getBucketGoing = false;
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

    changeType(item) {
        const oldItem = item;
        this.item = new BackupAccountCreateRequest();
        this.item.name = oldItem.name;
        this.item.type = oldItem.type;
        this.buckets = [];
    }

    getBuckets() {
        this.getBucketGoing = true;
        this.backupAccountService.listBuckets(this.item).subscribe(res => {
            this.buckets = res;
            this.getBucketGoing = false;
        }, error => {
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
            this.getBucketGoing = false;
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
