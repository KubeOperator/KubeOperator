import {Component, ElementRef, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {Terminal} from 'xterm';
import {Subscription} from 'rxjs';
import {LogService} from '../log.service';

@Component({
  selector: 'app-log-detail',
  templateUrl: './log-detail.component.html',
  styleUrls: ['./log-detail.component.css']
})
export class LogDetailComponent implements OnInit {
  @Input() showLogModal = false;
  @Output() showLogModalChange = new EventEmitter();
  @Input() logUrl = null;
  logText = '';

  constructor(private service: LogService) {
  }

  ngOnInit() {
  }

  loadLog() {
    if (this.logUrl != null) {
      this.service.getExecutionLog(this.logUrl).subscribe(data => {
        this.logText = data.data;
      });
    }
  }


  cancel() {
    this.showLogModal = false;
    this.showLogModalChange.emit(this.showLogModal);
  }

}
