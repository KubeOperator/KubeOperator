import {Component, OnInit, Output} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Execution} from './component/operater/execution';
import {LogService} from '../log/log.service';
import {DeployService} from './service/deploy.service';
import {Cluster} from '../cluster/cluster';
import {ClusterService} from '../cluster/cluster.service';
@Component({
  selector: 'app-deploy',
  templateUrl: './deploy.component.html',
  styleUrls: ['./deploy.component.css']
})
export class DeployComponent implements OnInit {

  currentCluster: Cluster;

  constructor(private route: ActivatedRoute, private clusterService: ClusterService, private executionService: LogService,
              private deployService: DeployService) {
  }


  ngOnInit() {
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
      // 更新cluster
      this.clusterService.getCluster(this.currentCluster.name).subscribe(cluster => {
        this.currentCluster = cluster;
        this.deployService.next(this.currentCluster.current_execution);
      });
    });
  }


}
