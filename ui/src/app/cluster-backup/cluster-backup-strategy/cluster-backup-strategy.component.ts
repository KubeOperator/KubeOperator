import {Component, OnInit} from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {ActivatedRoute} from '@angular/router';
import {ClusterBackupService} from '../cluster-backup.service';
import {BackupStrategy} from '../backup-strategy';
import {BackupStorageService} from '../../setting/backup-storage-setting/backup-storage.service';
import {BackupStorage} from '../../setting/backup-storage-setting/backup-storage';
import {NgForm} from '@angular/forms';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';


@Component({
  selector: 'app-cluster-backup-strategy',
  templateUrl: './cluster-backup-strategy.component.html',
  styleUrls: ['./cluster-backup-strategy.component.scss']
})
export class ClusterBackupStrategyComponent implements OnInit {

  constructor(private route: ActivatedRoute, private clusterBackupService: ClusterBackupService,
              private alertService: CommonAlertService, private backupStorageService: BackupStorageService) {
  }

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
    if (this.backupStrategy.id) {
      this.update();
    } else {
      this.create();
    }
  }

  getBackupStorage() {
    this.backupStorageService.listBackupStorage().subscribe(data => {
      this.loading = false;
      this.backupStorage = data;
    }, err => {
      this.loading = false;
    });
  }

  update() {
    if (!this.valid()) {
      return;
    }
    this.clusterBackupService.updateBackupStrategy(this.backupStrategy.project_id, this.backupStrategy).subscribe(data => {
      this.loading = false;
      this.alertService.showAlert('更新成功!', AlertLevels.SUCCESS);
      this.tipShow = false;
    }, err => {
      this.loading = false;
      this.alertService.showAlert('更新失败!' + err.reson + 'state code:' + err.status, AlertLevels.ERROR);
    });
  }

  create() {
    if (!this.valid()) {
      return;
    }
    this.clusterBackupService.createBackStrategy(this.backupStrategy).subscribe(data => {
      this.loading = false;
      this.alertService.showAlert('新增成功!', AlertLevels.SUCCESS);
      this.tipShow = false;
    }, err => {
      this.loading = false;
      this.alertService.showAlert('新增失败!' + err.reson + 'state code:' + err.status, AlertLevels.ERROR);
    });
  }

  valid() {
    if (this.backupStrategy.cron <= 0 || this.backupStrategy.cron > 300) {
      this.alertService.showAlert('备份间隔范围(1-300)', AlertLevels.ERROR);
      return false;
    } else if (this.backupStrategy.save_num <= 0 || this.backupStrategy.save_num > 100) {
      this.alertService.showAlert('保留份数范围(1-100)', AlertLevels.ERROR);
      return false;
    }
    return true;
  }

}
