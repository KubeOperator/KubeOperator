import {Component, OnInit} from '@angular/core';
import {CisService} from './cis.service';
import {Cluster} from '../cluster/cluster';
import {ClusterService} from '../cluster/cluster.service';
import {ActivatedRoute} from '@angular/router';
import {AlertLevels} from '../base/header/components/common-alert/alert';
import {CommonAlertService} from '../base/header/common-alert.service';

@Component({
  selector: 'app-cis',
  templateUrl: './cis.component.html',
  styleUrls: ['./cis.component.css']
})
export class CisComponent implements OnInit {

  loading = false;
  selectedItems = [];
  cises = [];
  currentCluster: Cluster;
  showDelete = false;
  resourceTypeName = 'CIS 扫描结果';

  constructor(private cisService: CisService, private clusterService: ClusterService,
              private route: ActivatedRoute, private alert: CommonAlertService) {
  }

  ngOnInit() {
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
      this.clusterService.getCluster(this.currentCluster.name).subscribe((d) => {
        this.currentCluster = d;
        this.listCis();
      });
    });
  }

  listCis() {
    this.loading = true;
    this.cisService.listCis(this.currentCluster.id).subscribe(data => {
      this.cises = data;
      this.loading = false;
    });
  }

  delete() {
    const promises: Promise<{}>[] = [];
    this.selectedItems.forEach(item => {
      promises.push(this.cisService.deleteCis(item.name).toPromise());
    });

    Promise.all(promises).then(data => {
      this.alert.showAlert('删除成功', AlertLevels.SUCCESS);
    }, res => {
      this.alert.showAlert('删除失败' + res.error.msg, AlertLevels.ERROR);
    }).finally(
      () => {
        this.showDelete = false;
        this.selectedItems = [];
        this.listCis();
      }
    );
  }

  check() {
    this.cisService.runCis(this.currentCluster.id).subscribe(res => {
      this.alert.showAlert(res.msg, AlertLevels.SUCCESS);
    });
  }

  download(id) {
    window.open('/api/v1/cluster/{id}/cisLog/download/'.replace('{id}', id));
  }
}
