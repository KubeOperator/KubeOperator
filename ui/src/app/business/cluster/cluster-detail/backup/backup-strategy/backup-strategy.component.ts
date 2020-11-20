import {Component, OnInit} from '@angular/core';
import {BackupFile, BackupStrategy} from '../cluster-backup';
import {BackupAccountService} from '../../../../setting/backup-account/backup-account.service';
import {BackupAccount} from '../../../../setting/backup-account/backup-account';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../../cluster';
import {BackupService} from '../backup.service';
import {CommonAlertService} from '../../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../../layout/common-alert/alert';
import {BackupFileService} from '../backup-file.service';
import {HttpClient} from '@angular/common/http';

@Component({
    selector: 'app-backup-strategy',
    templateUrl: './backup-strategy.component.html',
    styleUrls: ['./backup-strategy.component.css']
})
export class BackupStrategyComponent implements OnInit {

    backupStrategy: BackupStrategy = new BackupStrategy();
    backupAccounts: BackupAccount[] = [];
    currentCluster: Cluster;
    file;

    constructor(private backupAccountService: BackupAccountService,
                private route: ActivatedRoute,
                private backupService: BackupService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService,
                private backupFileService: BackupFileService) {
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.backupAccountService.listBy(this.currentCluster.projectName).subscribe(d => {
                this.backupAccounts = d.items;
                this.listBackupStrategy(this.currentCluster.name);
            });
        });
    }

    listBackupStrategy(clusterName) {
        this.backupService.getBy(clusterName).subscribe(s => {
            this.backupStrategy = s;
        });
    }

    upload(e): void {
        this.file = e.target.files[0];
    }

    onSubmit() {
        this.backupService.submit(this.backupStrategy).subscribe(res => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
            this.listBackupStrategy(this.currentCluster.name);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    onBackup() {
        const backupFile = new BackupFile();
        backupFile.clusterBackupStrategyID = this.backupStrategy.id;
        backupFile.clusterName = this.currentCluster.name;
        this.backupFileService.backup(backupFile).subscribe(res => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_BACKUP_START_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    onUploadFile() {
        const formData = new FormData();
        formData.append('file', this.file);
        formData.append('clusterName', this.currentCluster.name);
        this.backupFileService.localRestore(formData).subscribe(data => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_RESTORE_START_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
