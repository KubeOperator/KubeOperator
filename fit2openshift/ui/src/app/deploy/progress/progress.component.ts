import {Component, OnInit} from '@angular/core';
import {timer} from 'rxjs';
import {OperaterService} from '../../deploy/operater/operater.service';
import {LogService} from '../../log/log.service';

@Component({
  selector: 'app-progress',
  templateUrl: './progress.component.html',
  styleUrls: ['./progress.component.css']
})


export class ProgressComponent implements OnInit {
  currentProgress = 0;
  showProgress = false;

  constructor(private operaterService: OperaterService) {
  }

  ngOnInit() {
    this.operaterService.$executionQueue.subscribe(data => {
      this.showProgress = true;
      this.mock();
    });
  }

  mock() {
    const timers = setInterval(() => {
      this.currentProgress += 10;
      if (this.currentProgress === 100) {
        clearInterval(timers);
      }
    }, 2000);
  }

}
