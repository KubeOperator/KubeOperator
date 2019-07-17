import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, ActivatedRouteSnapshot} from '@angular/router';
import {Cluster} from '../cluster/cluster';
import {ClusterService} from '../cluster/cluster.service';

@Component({
  selector: 'app-overview',
  templateUrl: './overview.component.html',
  styleUrls: ['./overview.component.css']
})
export class OverviewComponent implements OnInit {

  currentCluster: Cluster;

  constructor(private route: ActivatedRoute, private clusterService: ClusterService) {
  }

  ngOnInit() {
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
      this.clusterService.getCluster(this.currentCluster.name).subscribe((d) => {
        this.currentCluster = d;
      });
    });
  }

}
