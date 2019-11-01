import {Component, OnInit} from '@angular/core';
import {SystemLogService} from '../system-log.service';
import {SystemLog} from '../system-log';

@Component({
  selector: 'app-system-log-list',
  templateUrl: './system-log-list.component.html',
  styleUrls: ['./system-log-list.component.css']
})
export class SystemLogListComponent implements OnInit {

  constructor(private systemLogService: SystemLogService) {
  }

  logs: SystemLog[] = [];
  currentPage: 1;
  totalItems: 0;
  keywords = '';
  level = 'all';

  ngOnInit() {
    this.refresh();
  }

  refresh() {
    const params = {
      level: this.level,
      currentPage: this.currentPage,
      keywords: this.keywords
    };
    this.systemLogService.searchLog(params).subscribe(data => {
      this.logs = data;
    });
  }

  onPageChange() {
    this.refresh();
  }

}
