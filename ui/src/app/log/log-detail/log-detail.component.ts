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
  @ViewChild('terminal', {static: true}) terminal: ElementRef;
  term: Terminal;

  constructor(private service: LogService) {
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
  }

  loadLog() {
    if (this.logUrl != null) {
      this.service.getExecutionLog(this.logUrl).subscribe(data => {
        this.term.write(data.data);
      });
    }
  }


  cancel() {
    this.showLogModal = false;
    this.showLogModalChange.emit(this.showLogModal);
  }

}
