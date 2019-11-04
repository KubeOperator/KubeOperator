import {Component, OnInit, ViewChild} from '@angular/core';
import {SystemLog} from '../system-log';

@Component({
  selector: 'app-system-log-detail',
  templateUrl: './system-log-detail.component.html',
  styleUrls: ['./system-log-detail.component.css']
})
export class SystemLogDetailComponent implements OnInit {

  constructor() {
  }

  open = false;
  log: SystemLog;
  ngOnInit() {
  }

  cancel() {
    this.open = false;
  }

}
