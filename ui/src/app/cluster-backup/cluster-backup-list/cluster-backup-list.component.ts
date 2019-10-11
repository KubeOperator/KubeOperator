import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {ClusterBackup} from '../cluster-backup';
import {ClusterBackupService} from '../cluster-backup.service';
import {ActivatedRoute} from '@angular/router';
import {TipService} from '../../tip/tip.service';
import {TipLevels} from '../../tip/tipLevels';
import {ConfirmAlertComponent} from '../../shared/common-component/confirm-alert/confirm-alert.component';


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


  constructor(private route: ActivatedRoute,  private clusterBackupService: ClusterBackupService,
               private tipService: TipService) {}

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
        this.tipService.showTip('删除成功', TipLevels.SUCCESS);
      }, error => {
        this.tipService.showTip('删除失败', TipLevels.ERROR);
      }).finally(
        () => {
          this.showDelete = false;
          this.selected = [];
          this.listClusterBackups();
        }
      );
      this.loading = false;
  }

  onRestore() {
      this.confirmAlert.setTitle('确认恢复');
      this.confirmAlert.setComment('确认以此备份恢复？');
      this.confirmAlert.opened = true;
  }

  restore() {
      this.clusterBackupService.restoreClusterBackup(this.selected[0]).subscribe(data => {
          this.tipService.showTip('恢复成功', TipLevels.SUCCESS);
      }, error1 => {

      });
      this.confirmAlert.close();
  }
}
