import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Cluster} from '../cluster';
import {ClusterService} from '../cluster.service';
import {Router} from '@angular/router';

@Component({
  selector: 'app-cluster-list',
  templateUrl: './cluster-list.component.html',
  styleUrls: ['./cluster-list.component.css']
})
export class ClusterListComponent implements OnInit {

  loading = true;
  clusters: Cluster[] = [];
  selectedRow: Cluster[] = [];
  @Output() addCluster = new EventEmitter<void>();

  constructor(private clusterService: ClusterService, private router: Router) {
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

  addNewCluster() {
    this.addCluster.emit();
  }

  goToLink(clusterId: string) {
    const linkUrl = ['fit2openshift', 'cluster', clusterId, 'overview'];
    this.router.navigate(linkUrl);
  }

}
