import {Component, OnInit} from '@angular/core';
import {ClusterStorageService} from './cluster-storage.service';
import {Cluster} from '../cluster/cluster';
import {ActivatedRoute} from '@angular/router';

@Component({
  selector: 'app-cluster-storage',
  templateUrl: './cluster-storage.component.html',
  styleUrls: ['./cluster-storage.component.css']
})
export class ClusterStorageComponent implements OnInit {
  currentCluster: Cluster;
  projectName = '';
  storages = [];
  loading = true;

  constructor(private clusterStorageService: ClusterStorageService, private route: ActivatedRoute,) {
  }

  ngOnInit() {
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
      this.projectName = this.currentCluster.name;
      this.listStorageClass();
    });
  }


  listStorageClass() {
    this.clusterStorageService.listClusterStorage(this.projectName).subscribe(res => {
      this.storages = res;
      this.loading = false;
    }, error1 => {
      this.loading = false;
    });
  }

}
