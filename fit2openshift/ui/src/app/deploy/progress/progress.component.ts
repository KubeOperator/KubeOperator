import {Component, Input, OnDestroy, OnInit} from '@angular/core';
import {Subscription, timer} from 'rxjs';
import {OperaterService} from '../operater/operater.service';
import {LogService} from '../../log/log.service';
import {WebsocketService} from '../term/websocket.service';
import {Execution} from '../operater/execution';
import {Cluster} from '../../cluster/cluster';
import {DeployService} from '../deploy.service';

export class ProgressMessage {
  id: string;
  progress: number;
  current_task: string;
  state: string;
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

  @Input() currentExecution: Execution;
  @Input() currentCluster: Cluster;


  constructor(private operaterService: OperaterService, private wsService: WebsocketService, private deployService: DeployService) {
  }

  ngOnInit() {
    this.deployService.$executionQueue.subscribe(data => {
      this.currentExecution = data;
      if (this.currentExecution === null) {
        this.showProgress = false;
      } else {
        // 判断是否完成
        this.showProgress = true;
        if (this.currentExecution.state !== 'SUCCESS' && this.currentExecution.state !== 'FAILURE') {
          this.subProgress();
        } else {
          this.fullProgress();
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


  fullProgress() {
    this.currentTask = this.currentExecution.current_task;
    this.currentProgress = this.currentExecution.progress * 100.0;
  }

  subProgress() {
    this.progressWsUrl = 'ws://' + window.location.host + this.currentExecution.progress_ws_url;
    this.progressSub = this.wsService.connect(this.progressWsUrl).subscribe(msg => {
      const m: ProgressMessage = JSON.parse(msg.data).message;
      this.currentProgress = m.progress * 100.0;
      this.currentTask = m.current_task;
      if (m.state !== 'SUCCESS' && m.state !== 'FAILURE') {
        this.deployService.nextState(false);
      } else {
        this.deployService.nextState(true);
      }
    });
  }


}
