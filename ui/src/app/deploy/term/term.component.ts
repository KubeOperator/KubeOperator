import {Component, ElementRef, Input, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {Terminal} from 'xterm';
import {Execution} from '../operater/execution';
import {WebsocketService} from './websocket.service';
import {OperaterService} from '../operater/operater.service';
import {Subject, Subscription} from 'rxjs';
import {Cluster} from '../../cluster/cluster';
import {LogService} from '../../log/log.service';
import {DeployService} from '../deploy.service';

@Component({
  selector: 'app-term',
  templateUrl: './term.component.html',
  styleUrls: ['./term.component.css']
})
export class TermComponent implements OnInit, OnDestroy {
  term: Terminal;
  logSub: Subscription;
  @ViewChild('terminal') terminal: ElementRef;
  currentExecution: Execution;
  @Input() currentCluster: Cluster;

  constructor(private wsService: WebsocketService, private operaterService: OperaterService,
              private executionService: LogService, private deployService: DeployService) {
  }

  ngOnInit() {
    this.deployService.$executionQueue.subscribe(data => {
      this.currentExecution = data;
      if (this.currentExecution === null) {
        this.resetTerm();
      } else {
        if (this.currentExecution.state !== 'SUCCESS' && this.currentExecution.state !== 'FAILURE') {
          this.subLog();
        } else {
          this.resetTerm();
        }
      }
      this.operaterService.$executionQueue.subscribe(e => {
        this.term.clear();
        if (this.logSub !== undefined && !this.logSub.closed) {
          this.logSub.unsubscribe();
        }
        this.currentExecution = e;
        this.subLog();
      });
    });

    this.term = new Terminal({
      cursorBlink: true,
      cols: 200,
      rows: 33,
      letterSpacing: 0,
      fontSize: 16
    });
    this.term.open(this.terminal.nativeElement);
  }


  resetTerm() {
    const banner = 'welcome to fit2openshift!';
    this.term.write(banner);
  }

  subLog() {
    this.term.clear();
    this.logSub = this.wsService.connect('ws://' + window.location.host + this.currentExecution.log_ws_url).subscribe(msg => {
      this.term.write(JSON.parse(msg.data).message);
    });
  }


  ngOnDestroy(): void {
    if (this.logSub !== undefined && !this.logSub.closed) {
      this.logSub.unsubscribe();
    }
  }

}
