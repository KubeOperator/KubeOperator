import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {BackupAccount} from '../backup-account';
import {BackupAccountService} from '../backup-account.service';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-backup-account-delete',
    templateUrl: './backup-account-delete.component.html',
    styleUrls: ['./backup-account-delete.component.css']
})
export class BackupAccountDeleteComponent extends BaseModelDirective<BackupAccount> implements OnInit {

    opened = false;
    items: BackupAccount[] = [];
    loading = false;
    @Output() deleted = new EventEmitter();

    constructor(private backupAccountService: BackupAccountService, private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private translateService: TranslateService) {
        super(backupAccountService);
    }

    ngOnInit(): void {
    }

    open(items) {
        this.items = items;
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
        this.items = [];
    }

    onSubmit() {
        this.loading = true;
        this.service.batch('delete', this.items).subscribe(data => {
            this.deleted.emit();
            this.opened = false;
            this.loading = false;
            this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.loading = false;
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
