import {Component, Input, OnInit} from '@angular/core';
import {Log} from '../log';
import {LogService} from '../log.service';
import {Cluster} from '../../cluster/cluster';
import {Execution} from '../../deploy/operater/execution';

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

  getColor(log: Log): string {
    const warn = '#FCFCAD';
    const error = '#FFAAAA';
    const info = '#AAFFCC';
    console.log(log.level);
    switch (log.level) {
      case 'INFO':
        return info;
      case 'WARN':
        return warn;
      case 'ERROR':
        return error;
    }

  }

  refresh() {
    this.listLog();
  }
}
