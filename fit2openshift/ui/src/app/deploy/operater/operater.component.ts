import {Component, Input, OnInit, Output} from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {OperaterService} from './operater.service';
import {Execution} from './execution';
import {LogService} from '../../log/log.service';
import {DeployService} from '../deploy.service';

@Component({
  selector: 'app-operater',
  templateUrl: './operater.component.html',
  styleUrls: ['./operater.component.css']
})
export class OperaterComponent implements OnInit {

  constructor(private operaterService: OperaterService, private logService: LogService, private deployService: DeployService) {
  }

  disableBtn = false;
  @Input() currentCluster: Cluster;
  currentExecution: Execution;

  ngOnInit() {
    this.deployService.$executionQueue.subscribe(data => {
      this.currentExecution = data;
      if (this.currentExecution !== null) {
        if (this.currentExecution.state !== 'SUCCESS' && this.currentExecution.state !== 'FAILURE') {
          this.disableBtn = true;
        }
      }

      this.deployService.$finished.subscribe(finished => {
        this.disableBtn = !finished;
      });
    });
  }

  install() {
    this.operaterService.executeOperate(this.currentCluster.name, 'install').subscribe(data => {
      this.currentExecution = data;
      this.operaterService.executionQueue.next(this.currentExecution);
    });
  }

  uninstall() {
    this.operaterService.executeOperate(this.currentCluster.name, 'uninstall').subscribe(data => {
      this.currentExecution = data;
      this.operaterService.executionQueue.next(this.currentExecution);
    });
  }
}
