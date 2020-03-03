import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {ActivatedRoute, Router} from '@angular/router';
import {ClusterBackupService} from '../cluster-backup.service';
import {BackupStrategy} from '../backup-strategy';
import {BackupStorageService} from '../../setting/backup-storage-setting/backup-storage.service';
import {BackupStorage} from '../../setting/backup-storage-setting/backup-storage';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {ConfirmAlertComponent} from '../../shared/common-component/confirm-alert/confirm-alert.component';
import {OperaterService} from '../../deploy/component/operater/operater.service';
import {ClusterHealthService} from '../../cluster-health/cluster-health.service';
import {SessionService} from '../../shared/session.service';


@Component({
  selector: 'app-cluster-backup-strategy',
  templateUrl: './cluster-backup-strategy.component.html',
  styleUrls: ['./cluster-backup-strategy.component.scss']
})
export class ClusterBackupStrategyComponent implements OnInit {

  constructor(private route: ActivatedRoute, private clusterBackupService: ClusterBackupService,
              private alertService: CommonAlertService, private operaterService: OperaterService,
              private backupStorageService: BackupStorageService, private router: Router,
              private clusterHealthService: ClusterHealthService, private sessionService: SessionService) {
  }

  tipShow = false;
  loading = false;
  backupStorage: BackupStorage[] = [];
  backupStrategy = new BackupStrategy();
  projectId = '';
  event: string = null;
  @ViewChild(ConfirmAlertComponent, {static: true}) confirmAlert: ConfirmAlertComponent;
  etcdHealth = true;
  permission;
  @Input() currentCluster: Cluster;


  ngOnInit() {
    // this.route.parent.data.subscribe(data => {
    //   this.currentCluster = data['cluster'];
    //   this.permission = this.sessionService.getItemPermission(this.currentCluster.item_name);
    //   this.projectId = this.currentCluster.id;
    //   this.getBackupStrategy();
    //   this.getBackupStorage();
    //   this.getClusterStatus();
    // });

    this.permission = this.sessionService.getItemPermission(this.currentCluster.item_name);
    this.projectId = this.currentCluster.id;
    this.getBackupStrategy();
    this.getBackupStorage();
    this.getClusterStatus();
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
    this.backupStorageService.listItemBackupStorage(this.currentCluster.item_name).subscribe(data => {
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
      this.backupStrategy.id = data.id;
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

  getClusterStatus() {
    if (this.currentCluster.status === 'READY') {
      this.etcdHealth = false;
      return;
    }
    this.clusterHealthService.listComponent(this.currentCluster.name).subscribe(res => {
      for (const ch of res) {
        if (ch.name.indexOf('etcd') !== -1 && ch.status !== 'RUNNING') {
          this.etcdHealth = false;
          break;
        }
      }
    });
  }

  onBackup() {
    if (!this.etcdHealth) {
      this.alertService.showAlert('集群ETCD不在运行状态 无法备份！', AlertLevels.ERROR);
      return;
    }
    if (this.backupStrategy.id == null) {
      this.alertService.showAlert('请先保存！', AlertLevels.ERROR);
      return;
    }
    this.confirmAlert.setTitle('确认备份');
    this.confirmAlert.setComment('立即开始备份？');
    this.confirmAlert.opened = true;
    this.event = 'backup';
  }

  handleBackup() {
    const params = {'backupStorageId': this.backupStrategy.backup_storage_id};
    this.handleEvent(params);
  }


  handleEvent(params?) {
    this.operaterService.executeOperate(this.currentCluster.name, this.event, params).subscribe(() => {
      this.redirect('deploy');
    });
    this.confirmAlert.close();
  }

  redirect(url: string) {
    if (url) {
      const linkUrl = ['cluster', this.currentCluster.name, url];
      this.router.navigate(linkUrl);
    }
  }

}
