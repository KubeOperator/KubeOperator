import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {ClusterBackup} from '../cluster-backup';
import {ClusterBackupService} from '../cluster-backup.service';
import {ActivatedRoute} from '@angular/router';
import {ConfirmAlertComponent} from '../../shared/common-component/confirm-alert/confirm-alert.component';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';


@Component({
  selector: 'app-cluster-backup-list',
  templateUrl: './cluster-backup-list.component.html',
  styleUrls: ['./cluster-backup-list.component.scss']
})
export class ClusterBackupListComponent implements OnInit {

  @Input() currentCluster: Cluster;
  loading = true;
  showDelete = false;
  items: ClusterBackup[] = [];
  selected: ClusterBackup[] = [];
  resourceTypeName = '备份';
  projectId = '';
  @ViewChild(ConfirmAlertComponent, {static: true}) confirmAlert: ConfirmAlertComponent;


  constructor(private route: ActivatedRoute, private clusterBackupService: ClusterBackupService,
              private alertService: CommonAlertService) {
  }

  ngOnInit() {
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
      this.projectId = this.currentCluster.id;
      this.listClusterBackups();
    });
  }

  listClusterBackups() {
    this.clusterBackupService.listClusterBackup(this.projectId).subscribe(data => {
      this.items = data;
    }, error1 => {

    });
  }

  delete() {
    const promises: Promise<{}>[] = [];
    this.loading = true;
    this.selected.forEach(item => {
      promises.push(this.clusterBackupService.deleteClusterBackup(item.id).toPromise());
    });

    Promise.all(promises).then(data => {
      this.alertService.showAlert('删除成功', AlertLevels.SUCCESS);
    }, error => {
      this.alertService.showAlert('删除失败', AlertLevels.ERROR);
    }).finally(
      () => {
        this.showDelete = false;
        this.selected = [];
        this.listClusterBackups();
      }
    );
    this.loading = false;
  }

  restore() {
    this.confirmAlert.setTitle('确认恢复');
    this.confirmAlert.setComment('确认以此备份恢复？');
    this.clusterBackupService.restoreClusterBackup(this.selected[0]).subscribe(data => {
      this.alertService.showAlert('恢复成功', AlertLevels.SUCCESS);
    }, error1 => {

    });
  }
}
