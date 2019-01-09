import {Component, Input, OnInit} from '@angular/core';
import {Log} from '../log';
import {LogService} from '../log.service';

@Component({
  selector: 'app-log-list',
  templateUrl: './log-list.component.html',
  styleUrls: ['./log-list.component.css']
})
export class LogListComponent implements OnInit {

  loading = true;
  logs: Log[] = [];
  @Input() clusterId: string;

  constructor(private logService: LogService) {
  }

  ngOnInit() {
    this.listLog();

    this.logService.messages.subscribe(data => {
      console.log(data.data);
    });

  }

  listLog() {
    this.loading = true;
    this.logService.getLogs(this.clusterId).subscribe(data => {
      this.loading = false;
      this.logs = data;
    }, error => {
      this.loading = false;
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
