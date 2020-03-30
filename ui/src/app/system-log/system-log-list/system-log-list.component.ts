import {Component, OnInit, ViewChild} from '@angular/core';
import {SystemLogService} from '../system-log.service';
import {SystemLog} from '../system-log';
import {SystemLogDetailComponent} from '../system-log-detail/system-log-detail.component';

@Component({
  selector: 'app-system-log-list',
  templateUrl: './system-log-list.component.html',
  styleUrls: ['./system-log-list.component.css']
})
export class SystemLogListComponent implements OnInit {

  constructor(private systemLogService: SystemLogService) {
  }

  logs: SystemLog[] = [];
  currentPage = 1;
  totalItems: number;
  keywords = '';
  level = 'INFO';
  loading = true;
  limit = '1';
  size = 10;
  @ViewChild(SystemLogDetailComponent, {static: true})
  detail: SystemLogDetailComponent;

  ngOnInit() {
    this.refresh();
  }

  refresh(rest_page?: boolean) {
    if (rest_page) {
      this.currentPage = 1;
    }
    const params = {
        level: this.level,
        page: this.currentPage,
        size: this.size,
        keywords: this.keywords,
        limit: this.limit,
      }
    ;
    this.loading = true;
    this.systemLogService.searchLog(params).subscribe(data => {
      this.logs = data.items;
      this.totalItems = data.total;
      this.loading = false;
    });
  }

  showDetail(log: SystemLog) {
    this.detail.log = log;
    this.detail.open = true;
  }

  onPageChange() {
    this.refresh();
  }
}
