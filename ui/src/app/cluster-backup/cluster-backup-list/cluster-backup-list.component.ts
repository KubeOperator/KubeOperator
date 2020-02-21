import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {ClusterBackup} from '../cluster-backup';
import {ClusterBackupService} from '../cluster-backup.service';
import {ActivatedRoute, Router} from '@angular/router';
import {ConfirmAlertComponent} from '../../shared/common-component/confirm-alert/confirm-alert.component';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {OperaterService} from '../../deploy/component/operater/operater.service';


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
  baseRoute: string;


  constructor(private route: ActivatedRoute, private clusterBackupService: ClusterBackupService,
              private alertService: CommonAlertService, private operaterService: OperaterService, private router: Router) {
  }

  ngOnInit() {
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
      this.projectId = this.currentCluster.id;
      this.baseRoute = 'item/' + this.currentCluster.item_name + '/cluster/' + this.currentCluster.name;
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
      this.router.navigate([this.baseRoute + '/deploy']);
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
