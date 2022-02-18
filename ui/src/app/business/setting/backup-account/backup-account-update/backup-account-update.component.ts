import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {BackupAccount, BackupAccountUpdateRequest} from '../backup-account';
import {BackupAccountService} from '../backup-account.service';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-backup-account-update',
    templateUrl: './backup-account-update.component.html',
    styleUrls: ['./backup-account-update.component.css']
})
export class BackupAccountUpdateComponent extends BaseModelDirective<BackupAccount> implements OnInit {

    opened = false;
    item: BackupAccountUpdateRequest = new BackupAccountUpdateRequest();
    buckets: [] = [];
    isSubmitGoing = false;
    @Output() updated = new EventEmitter();
    loadingBucket = false;

    constructor(private backupAccountService: BackupAccountService,
                private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
        super(backupAccountService);
    }

    ngOnInit(): void {
    }

    open(item) {
        Object.assign(this.item, item);
        if (item.type !== 'SFTP') {
            this.item.bucket = '';
        }
        this.opened = true;
    }

    getBuckets() {
        this.loadingBucket = true;
        this.backupAccountService.listBuckets(this.item).subscribe(res => {
            this.loadingBucket = false;
            this.buckets = res;
        }, error => {
            this.loadingBucket = false;
            this.buckets = [];
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    onCancel() {
        this.opened = false;
        this.item = new BackupAccountUpdateRequest();
    }

    onSubmit() {
        this.isSubmitGoing = true;
        this.backupAccountService.update(this.item.name, this.item).subscribe(res => {
            this.isSubmitGoing = false;
            this.updated.emit();
            this.onCancel();
            this.commonAlertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.msg ? error.msg : error.error.msg, AlertLevels.ERROR);
        });
    }
}
