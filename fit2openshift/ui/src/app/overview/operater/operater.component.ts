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

  ngOnInit() {
    if (this.currentCluster.current_task_id !== '') {
      this.getCurrentExecution();
    }
  }

  getCurrentExecution() {
    this.logService.getExecution(this.currentCluster.name, this.currentCluster.current_task_id).subscribe(data => {
      this.currentExecution = data;
      if (this.currentExecution) {
        this.operaterService.executionQueue.next(this.currentExecution);
      }
    });
  }

  startDeploy() {
    this.operaterService.startDeploy(this.currentCluster.name).subscribe(data => {
      this.currentExecution = data;
    });
  }
}
