import {Component, Input, OnInit, Output} from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {OperaterService} from './operater.service';
import {Execution} from './execution';
import {LogService} from '../../log/log.service';

@Component({
  selector: 'app-operater',
  templateUrl: './operater.component.html',
  styleUrls: ['./operater.component.css']
})
export class OperaterComponent implements OnInit {

  constructor(private operaterService: OperaterService, private logService: LogService) {
  }

  @Input() currentCluster: Cluster;
  @Output() currentExecution: Execution;
  clusterStatus = 'PENDING';

  ngOnInit() {
    this.getClusterStatus();
  }

  getClusterStatus() {
    this.logService.listExecutions(this.currentCluster.name).subscribe(data => {
      if (data.length > 0) {
        this.currentExecution = data[0];
        this.operaterService.executionQueue.next(this.currentExecution);
        this.clusterStatus = this.currentExecution.state;
      }
    });
  }

  startDeploy() {
    this.operaterService.startDeploy(this.currentCluster.name).subscribe(data => {
      this.currentExecution = data;
      this.clusterStatus = this.currentExecution.state;
    });
  }
}
