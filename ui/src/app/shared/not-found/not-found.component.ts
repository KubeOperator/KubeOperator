import {Component, OnDestroy, OnInit} from '@angular/core';
import {Router} from '@angular/router';
import {CommonRoutes} from '../shared.const';

const defaultInterval = 1000;
const defaultLeftTime = 5;

@Component({
  selector: 'app-not-found',
  templateUrl: './not-found.component.html',
  styleUrls: ['./not-found.component.css']
})
export class NotFoundComponent implements OnInit, OnDestroy {

  leftSeconds: number = defaultLeftTime;
  timeInterval: any = null;

  constructor(private router: Router) {
  }

  ngOnInit() {
    if (!this.timeInterval) {
      this.timeInterval = setInterval(interval => {
        this.leftSeconds--;
        if (this.leftSeconds <= 0) {
          this.router.navigateByUrl(CommonRoutes.F2O_DEFAULT  );
          clearInterval(this.timeInterval);
        }
      }, defaultInterval);
    }
  }

  ngOnDestroy(): void {
    if (this.timeInterval) {
      clearInterval(this.timeInterval);
    }
  }


}
