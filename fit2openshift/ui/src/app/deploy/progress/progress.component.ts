import {Component, Input, OnDestroy, OnInit} from '@angular/core';
import {Subscription, timer} from 'rxjs';
import {OperaterService} from '../operater/operater.service';
import {WebsocketService} from '../term/websocket.service';
import {Execution} from '../operater/execution';
import {Cluster} from '../../cluster/cluster';
import {DeployService} from '../deploy.service';

export class ProgressMessage {
  id: string;
  progress: number;
  current_task: string;
  state: string;
  operation: string;
}


@Component({
  selector: 'app-progress',
  templateUrl: './progress.component.html',
  styleUrls: ['./progress.component.css']
})


export class ProgressComponent implements OnInit, OnDestroy {
  currentProgress = 0;
  showProgress = true;
  progressWsUrl: string;
  progressSub: Subscription;
  currentTask: string;
  currentDeploy: string;
  showLoading = false;

  @Input() currentExecution: Execution;
  @Input() currentCluster: Cluster;


  constructor(private operaterService: OperaterService, private wsService: WebsocketService, private deployService: DeployService) {
  }


  ngOnInit() {
    this.deployService.$executionQueue.subscribe(data => {
      this.currentExecution = data;
      if (this.currentExecution === null) {
        this.currentDeploy = this.getDeploymentName('none');
      } else {
        // 判断是否完成
        if (this.currentExecution.state !== 'SUCCESS' && this.currentExecution.state !== 'FAILURE') {
          this.subProgress();
          this.currentDeploy = this.getDeploymentName(this.currentExecution.state);
        } else {
          this.currentDeploy = this.getDeploymentName('none');
        }
      }
      this.operaterService.$executionQueue.subscribe(d => {
        if (this.progressSub && !this.progressSub.closed) {
          this.progressSub.unsubscribe();
        }
        this.currentExecution = d;
        this.subProgress();
      });
    });
  }

  ngOnDestroy(): void {
    if (this.progressSub !== undefined && !this.progressSub.closed) {
      this.progressSub.unsubscribe();
    }
  }

  subProgress() {
    this.showLoading = true;
    this.progressWsUrl = 'ws://' + window.location.host + this.currentExecution.progress_ws_url;
    this.progressSub = this.wsService.connect(this.progressWsUrl).subscribe(msg => {
      const m: ProgressMessage = JSON.parse(msg.data).message;
      this.currentProgress = m.progress * 100.0;
      this.currentTask = m.current_task;
      this.currentDeploy = this.getDeploymentName(m.operation);
      if (m.state !== 'SUCCESS' && m.state !== 'FAILURE') {
        this.deployService.nextState(false);
      } else {
        this.deployService.nextState(true);
        this.currentDeploy = this.getDeploymentName('none');
        this.showLoading = false;
      }
    });
  }

  getDeploymentName(str: string) {
    switch (str) {
      case 'install':
        return '部署OKD集群';
      case 'change':
        return '变更节点';
      case'uninstall':
        return '卸载OKD集群';
      default:
        return '无';
    }
  }


}
