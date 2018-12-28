import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Cluster} from '../cluster';
import {ClusterService} from '../cluster.service';
import {Router} from '@angular/router';
import {TipService} from '../../tip/tip.service';
import {TipLevels} from '../../tip/tipLevels';

@Component({
  selector: 'app-cluster-list',
  templateUrl: './cluster-list.component.html',
  styleUrls: ['./cluster-list.component.css']
})
export class ClusterListComponent implements OnInit {

  loading = true;
  clusters: Cluster[] = [];
  selected: Cluster[] = [];
  @Output() addCluster = new EventEmitter<void>();

  constructor(private clusterService: ClusterService, private router: Router, private tipService: TipService) {
  }

  ngOnInit() {
    this.listCluster();
  }

  listCluster() {
    this.clusterService.listCluster().subscribe(data => {
      this.clusters = data;
      this.loading = false;
    }, error => {
      this.loading = false;
    });
  }

  deleteClusters() {
    if (!(this.selected.length > 0)) {
      this.tipService.showTip('请选择要删除的集群!', TipLevels.ERROR);
      return;
    }
    this.loading = true;
    this.selected.forEach(cluster => {
      let flag = false;
      this.clusterService.deleteCluster(cluster.name).subscribe(data => {
        flag = true;
      });
      while (!flag) {

      }
    });
    this.loading = false;
  }

  addNewCluster() {
    this.addCluster.emit();
  }

  goToLink(clusterId: string) {
    const linkUrl = ['fit2openshift', 'cluster', clusterId, 'overview'];
    this.router.navigate(linkUrl);
  }

}
