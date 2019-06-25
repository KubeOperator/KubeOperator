import {Component, Input, OnDestroy, OnInit} from '@angular/core';
import {Subscription, timer} from 'rxjs';
import {OperaterService} from '../operater/operater.service';
import {WebsocketService} from '../term/websocket.service';
import {Execution} from '../operater/execution';
import {DeployService} from '../../service/deploy.service';
import {DeplayUtilService} from '../../service/deplay-util.service';
import {Cluster} from '../../../cluster/cluster';

export class ProgressMessage {
  progress: number;
  current_play: string;
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


  constructor(private operaterService: OperaterService, private wsService: WebsocketService,
              private deployService: DeployService, private deplayUtil: DeplayUtilService) {
  }


  ngOnInit() {
    this.deployService.$executionQueue.subscribe(data => {
      this.currentExecution = data;
      if (!this.currentExecution) {
        this.currentDeploy = '暂无';
      } else {
        // 判断是否完成
        if (!this.deplayUtil.execution_is_complated(this.currentExecution.state)) {
          this.subProgress();
          this.currentDeploy = this.currentExecution.operation;
        } else {
          this.currentDeploy = '暂无';
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
      const m: ProgressMessage = JSON.parse(JSON.parse(msg.data).message);
      this.currentProgress = m.progress;
      this.currentTask = m.current_play;
      this.currentDeploy = m.operation;
      if (!this.deplayUtil.execution_is_complated(m.state)) {
        this.deployService.nextState(false);
      } else {
        this.deployService.nextState(true);
        this.currentDeploy = m.operation;
        this.showLoading = false;
      }
    });
  }


}
