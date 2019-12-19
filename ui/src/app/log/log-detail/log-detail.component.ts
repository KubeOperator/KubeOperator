import {Component, ElementRef, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {Terminal} from 'xterm';
import {Subscription} from 'rxjs';
import {LogService} from '../log.service';
import {logging} from 'selenium-webdriver';

@Component({
  selector: 'app-log-detail',
  templateUrl: './log-detail.component.html',
  styleUrls: ['./log-detail.component.css']
})
export class LogDetailComponent implements OnInit {
  @Input() showLogModal = false;
  @Output() showLogModalChange = new EventEmitter();
  @Input() logUrl = null;
  @ViewChild('terminal', {static: true}) terminal: ElementRef;
  loading = true;
  term: Terminal;

  constructor(private service: LogService) {
  }

  ngOnInit() {
    this.term = new Terminal({
      cursorBlink: false,
      disableStdin: true,
      cursorStyle: 'bar',
      cols: 110,
      rows: 25,
      letterSpacing: 1,
      scrollback: 9999999,
    });
    this.term.open(this.terminal.nativeElement);
  }

  loadLog() {
    this.showLogModal = true;
    this.loading = true;
    if (this.logUrl != null) {
      this.term.reset();
      this.service.getExecutionLog(this.logUrl).subscribe(data => {
        this.loading = false;
        this.term.write(data.data);
      });
    }
  }


  cancel() {
    this.showLogModal = false;
    this.showLogModalChange.emit(this.showLogModal);
  }

}
