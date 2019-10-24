import {Component, ElementRef, Input, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {Terminal} from 'xterm';
import {Execution} from '../operater/execution';
import {WebsocketService} from './websocket.service';
import {OperaterService} from '../operater/operater.service';
import {Subject, Subscription} from 'rxjs';
import {LogService} from '../../../log/log.service';
import {DeployService} from '../../service/deploy.service';
import {Cluster} from '../../../cluster/cluster';
import {DeplayUtilService} from '../../service/deplay-util.service';

@Component({
  selector: 'app-term',
  templateUrl: './term.component.html',
  styleUrls: ['./term.component.css']
})
export class TermComponent implements OnInit, OnDestroy {
  term: Terminal;
  logSub: Subscription;
  @ViewChild('terminal', {static: true}) terminal: ElementRef;
  currentExecution: Execution;
  @Input() currentCluster: Cluster;

  constructor(private wsService: WebsocketService, private operaterService: OperaterService,
              private executionService: LogService, private deployService: DeployService, private deployUtil: DeplayUtilService) {
  }

  ngOnInit() {
    this.term = new Terminal({
      cursorBlink: false,
      disableStdin: true,
      cursorStyle: 'bar',
      cols: 120,
      rows: 30,
      letterSpacing: 1,
      fontSize: 16
    });
    this.term.open(this.terminal.nativeElement);
    this.resetTerm();
    this.deployService.$executionQueue.subscribe(data => {
      this.currentExecution = data;
      if (!this.deployUtil.execution_is_complated(this.currentExecution.state)) {
        this.subLog();
      }
    });
  }


  resetTerm() {
    const banner = 'Welcome to KubeOperator';
    this.term.write(banner);
  }

  subLog() {
    this.term.clear();
    this.logSub = this.wsService.connect('ws://' + window.location.host + this.currentExecution.log_ws_url).subscribe(msg => {
      this.term.write(JSON.parse(msg.data).message);
    });
  }


  ngOnDestroy(): void {
    if (this.logSub) {
      this.logSub.unsubscribe();
    }
  }
}
