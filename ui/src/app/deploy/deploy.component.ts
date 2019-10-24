import {Component, OnInit, Output} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {LogService} from '../log/log.service';
import {DeployService} from './service/deploy.service';
import {Cluster} from '../cluster/cluster';
import {ClusterService} from '../cluster/cluster.service';
import {DeplayUtilService} from './service/deplay-util.service';

@Component({
  selector: 'app-deploy',
  templateUrl: './deploy.component.html',
  styleUrls: ['./deploy.component.css']
})
export class DeployComponent implements OnInit {

  currentCluster: Cluster;
  placeholder = false;

  constructor(private route: ActivatedRoute, private clusterService: ClusterService, private executionService: LogService,
              private deployService: DeployService, private deployUtil: DeplayUtilService) {
  }

  ngOnInit() {
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
      // 更新cluster
      this.clusterService.getCluster(this.currentCluster.name).subscribe(cluster => {
        this.currentCluster = cluster;
        if (this.currentCluster.current_execution && !this.deployUtil.execution_is_complated(this.currentCluster.current_execution.state)) {
          this.deployService.next(this.currentCluster.current_execution);
          this.placeholder = false;
        } else {
          this.placeholder = true;
        }
      });
    });
  }


}
