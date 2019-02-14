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
  deleteModal = false;
  selectedCluster: Cluster = new Cluster();

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


  deleteCluster(cluster: Cluster) {
    this.deleteModal = true;
    this.selectedCluster = cluster;
  }

  confirmDelete() {
    this.clusterService.deleteCluster(this.selectedCluster.name).subscribe(data => {
      this.deleteModal = false;
      this.listCluster();
      this.tipService.showTip('删除成功！', TipLevels.SUCCESS);
    }, error => {
      this.deleteModal = false;
      this.tipService.showTip('删除失败！ msg:' + error, TipLevels.ERROR);
    });
  }


  addNewCluster() {
    this.addCluster.emit();
  }

  goToLink(clusterName: string) {
    const linkUrl = ['fit2openshift', 'cluster', clusterName, 'overview'];
    this.router.navigate(linkUrl);
  }

}
