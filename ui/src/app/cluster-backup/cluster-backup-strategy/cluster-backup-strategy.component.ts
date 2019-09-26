import { Component, OnInit } from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {ActivatedRoute} from '@angular/router';
import {ClusterBackupService} from '../cluster-backup.service';
import {BackupStrategy} from '../backup-strategy';
import {TipService} from '../../tip/tip.service';
import {TipLevels} from '../../tip/tipLevels';
import {BackupStorageService} from '../../setting/backup-storage-setting/backup-storage.service';
import {BackupStorage} from '../../setting/backup-storage-setting/backup-storage';


@Component({
  selector: 'app-cluster-backup-strategy',
  templateUrl: './cluster-backup-strategy.component.html',
  styleUrls: ['./cluster-backup-strategy.component.scss']
})
export class ClusterBackupStrategyComponent implements OnInit {

  constructor(private route: ActivatedRoute,  private clusterBackupService: ClusterBackupService,
               private tipService: TipService, private backupStorageService: BackupStorageService) {}
  tipShow = false;
  loading = false;
  currentCluster: Cluster;
  backupStorage: BackupStorage[] = [];
  backupStrategy = new BackupStrategy();
  projectId = '';

  ngOnInit() {
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
      this.projectId = this.currentCluster.id;
      this.getBackupStrategy();
      this.getBackupStorage();
    });
  }

  getBackupStrategy() {
      this.clusterBackupService.listBackupStrategy(this.projectId).subscribe(data => {
          this.backupStrategy = data;
      }, error => {
          this.backupStrategy = new BackupStrategy();
          this.backupStrategy.project_id = this.projectId;
      });
  }

  onCommit() {
      this.clusterBackupService.createBackStrategy(this.backupStrategy).subscribe(data => {
        this.loading = false;
        this.tipService.showTip('新增成功!', TipLevels.SUCCESS);
        this.tipShow = false;
      }, err => {
        this.loading = false;
        this.tipService.showTip('新增失败!' + err.reson + 'state code:' + err.status, TipLevels.ERROR);
      });
  }

  getBackupStorage() {
      this.backupStorageService.listBackupStorage().subscribe(data => {
        this.loading = false;
        this.backupStorage = data;
      }, err => {
        this.loading = false;
      });
  }

}
