import {Component, OnInit} from '@angular/core';
import {BaseModelDirective} from '../../../../../shared/class/BaseModelDirective';
import {BackupFile} from '../cluster-backup';
import {ActivatedRoute} from '@angular/router';
import {CommonAlertService} from '../../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {BackupFileService} from '../backup-file.service';
import {Cluster} from '../../../cluster';
import {AlertLevels} from '../../../../../layout/common-alert/alert';

@Component({
    selector: 'app-backup-list',
    templateUrl: './backup-list.component.html',
    styleUrls: ['./backup-list.component.css']
})
export class BackupListComponent extends BaseModelDirective<BackupFile> implements OnInit {

    currentCluster: Cluster;
    items: BackupFile[] = [];
    deleteOpen = false;
    deleteName = '';
    loading = false;

    constructor(private route: ActivatedRoute,
                private backupFileService: BackupFileService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
        super(backupFileService);
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.pageby();
        });
    }

    pageby() {
        this.backupFileService.pageBy(this.page, this.size, this.currentCluster.name).subscribe(d => {
            this.items = d.items;
        });
    }

    restore(name) {
        const restoreRequest = new BackupFile();
        restoreRequest.clusterName = this.currentCluster.name;
        restoreRequest.name = name;
        this.backupFileService.restore(restoreRequest).subscribe(res => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_RESTORE_START_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    deleteFile(name) {
        this.deleteName = name;
        this.deleteOpen = true;
    }

    cancelDelete() {
        this.deleteName = '';
        this.deleteOpen = false;
    }

    submitDelete() {
        this.loading = true;
        this.backupFileService.delete(this.deleteName).subscribe(res => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);
            this.loading = false;
            this.pageby();
            this.cancelDelete();
        }, error => {
            this.loading = false;
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
            this.cancelDelete();
        });
    }
}
