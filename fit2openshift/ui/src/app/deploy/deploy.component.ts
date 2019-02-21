import {Component, OnInit, Output} from '@angular/core';
import {Cluster} from '../cluster/cluster';
import {ActivatedRoute} from '@angular/router';
import {ClusterService} from '../cluster/cluster.service';
import {Execution} from './operater/execution';
import {LogService} from '../log/log.service';
import {Subject} from 'rxjs';
import {Tip} from '../tip/tip';
import {DeployService} from './deploy.service';

@Component({
  selector: 'app-deploy',
  templateUrl: './deploy.component.html',
  styleUrls: ['./deploy.component.css']
})
export class DeployComponent implements OnInit {

  currentCluster: Cluster;
  currentExecution: Execution;

  constructor(private route: ActivatedRoute, private clusterService: ClusterService, private executionService: LogService, private deployService: DeployService) {
  }


  ngOnInit() {
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
      this.currentExecution = null;
      // 更新cluster
      this.clusterService.getCluster(this.currentCluster.name).subscribe(cluster => {
        this.currentCluster = cluster;
        if (cluster.current_task_id != null) {
          this.executionService.getExecution(this.currentCluster.name, this.currentCluster.current_task_id).subscribe(execution => {
            this.deployService.next(execution);
          });
        } else {
          this.deployService.next(null);
        }
      });
    });
  }


}
