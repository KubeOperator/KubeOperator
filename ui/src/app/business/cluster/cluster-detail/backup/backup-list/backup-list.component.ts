import {Component, OnInit} from '@angular/core';
import {BaseModelComponent} from '../../../../../shared/class/BaseModelComponent';
import {BackupFile} from '../cluster-backup';
import {ActivatedRoute} from '@angular/router';
import {CommonAlertService} from '../../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {BackupFileService} from '../backup-file.service';
import {Cluster} from '../../../cluster';

@Component({
    selector: 'app-backup-list',
    templateUrl: './backup-list.component.html',
    styleUrls: ['./backup-list.component.css']
})
export class BackupListComponent extends BaseModelComponent<BackupFile> implements OnInit {

    currentCluster: Cluster;
    items: BackupFile[] = [];

    constructor(private route: ActivatedRoute,
                private backupFileService: BackupFileService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
        super(backupFileService);
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.backupFileService.pageBy(this.page, this.size, this.currentCluster.name).subscribe(d => {
                this.items = d.items;
            });
        });
    }

}
