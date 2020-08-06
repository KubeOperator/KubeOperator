import {Component, OnInit} from '@angular/core';
import {BackupStrategy} from '../cluster-backup';
import {BackupAccountService} from '../../../../setting/backup-account/backup-account.service';
import {BackupAccount} from '../../../../setting/backup-account/backup-account';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../../cluster';
import {BackupService} from '../backup.service';
import {CommonAlertService} from '../../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../../layout/common-alert/alert';

@Component({
    selector: 'app-backup-strategy',
    templateUrl: './backup-strategy.component.html',
    styleUrls: ['./backup-strategy.component.css']
})
export class BackupStrategyComponent implements OnInit {

    backupStrategy: BackupStrategy = new BackupStrategy();
    backupAccounts: BackupAccount[] = [];
    currentCluster: Cluster;

    constructor(private backupAccountService: BackupAccountService,
                private route: ActivatedRoute,
                private backupService: BackupService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {

    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.backupAccountService.listBy(this.currentCluster.projectName).subscribe(d => {
                this.backupAccounts = d.items;
                this.backupService.getBy(this.currentCluster.name).subscribe(s => {
                    this.backupStrategy = s;
                });
            });
        });
    }


    onSubmit() {
        this.backupService.submit(this.backupStrategy).subscribe(res => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
