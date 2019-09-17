import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {Log} from '../log';
import {LogService} from '../log.service';
import {Cluster} from '../../cluster/cluster';
import {HostListComponent} from '../../host/host-list/host-list.component';
import {LogDetailComponent} from '../log-detail/log-detail.component';
import {Execution} from '../../deploy/component/operater/execution';

@Component({
  selector: 'app-log-list',
  templateUrl: './log-list.component.html',
  styleUrls: ['./log-list.component.css']
})
export class LogListComponent implements OnInit {

  loading = true;
  logs: Log[] = [];
  executions: Execution[] = [];
  @Input() currentCluster: Cluster;
  showDetail = false;
  @ViewChild(LogDetailComponent, { static: true })
  child: LogDetailComponent;

  constructor(private logService: LogService) {
  }

  ngOnInit() {
    this.listLog();
  }

  listLog() {
    this.loading = true;
    this.logService.listExecutions(this.currentCluster.name).subscribe(data => {
      this.loading = false;
      this.executions = data;
    });
  }

  showLogDetail(execution: Execution) {
    this.showDetail = true;
    this.child.logUrl = execution.log_url;
    this.child.loadLog();
  }

  refresh() {
    this.listLog();
  }
}
