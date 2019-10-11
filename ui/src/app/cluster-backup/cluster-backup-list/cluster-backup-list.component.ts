import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {ClusterBackup} from '../cluster-backup';
import {ClusterBackupService} from '../cluster-backup.service';
import {ActivatedRoute} from '@angular/router';
import {TipService} from '../../tip/tip.service';
import {TipLevels} from '../../tip/tipLevels';
import {ConfirmAlertComponent} from '../../shared/common-component/confirm-alert/confirm-alert.component';
import {OperaterService} from '../../deploy/component/operater/operater.service';
import {Router} from '@angular/router';

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
  event: string = null;
  @ViewChild(ConfirmAlertComponent, {static: true}) confirmAlert: ConfirmAlertComponent;


  constructor(private route: ActivatedRoute,  private clusterBackupService: ClusterBackupService,
               private tipService: TipService, private operaterService: OperaterService, private router: Router) {}

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
      this.event = 'restore';
  }

  handleRestore() {
      const params = {'clusterBackupId': this.selected[0].id};
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
      const linkUrl = ['kubeOperator', 'cluster', this.currentCluster.name, url];
      this.router.navigate(linkUrl);
    }
  }
}
