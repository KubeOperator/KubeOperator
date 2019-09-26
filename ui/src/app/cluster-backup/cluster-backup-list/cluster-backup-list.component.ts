import {Component, Input, OnInit} from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {ClusterBackup} from '../cluster-backup';
import {ClusterBackupService} from '../cluster-backup.service';
import {ActivatedRoute} from '@angular/router';
import {TipService} from '../../tip/tip.service';
import {TipLevels} from "../../tip/tipLevels";

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
      this.selected.forEach(item => {
          promises.push(this.clusterBackupService.deleteClusterBackup(item.name).toPromise());
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
  }
}
