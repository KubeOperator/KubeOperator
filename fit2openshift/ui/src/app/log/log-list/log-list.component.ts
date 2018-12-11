import {Component, Input, OnInit} from '@angular/core';
import {Log} from '../log';

@Component({
  selector: 'app-log-list',
  templateUrl: './log-list.component.html',
  styleUrls: ['./log-list.component.css']
})
export class LogListComponent implements OnInit {

  loading = true;
  logs: Log[] = [];
  @Input() clusterId: string;

  constructor() {
  }

  ngOnInit() {
    this.listLog();
  }

  listLog() {
    const log: Log = new Log();
    log.date = '2018-12-11 09:33:33';
    log.level = 'INFO';
    log.message = 'create a job ......';
    this.logs.push(log);
    this.loading = false;
  }

}
