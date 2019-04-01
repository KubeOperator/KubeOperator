import {Component, OnDestroy, OnInit} from '@angular/core';
import {Tip} from './tip';
import {TipService} from './tip.service';
import {TipLevels} from './tipLevels';

const defaultInterval = 1000;
const defaultLeftTime = 5;

@Component({
  selector: 'app-tip',
  templateUrl: './tip.component.html',
  styleUrls: ['./tip.component.css']
})
export class TipComponent implements OnInit {
  currentTip: Tip;
  currentTipLevel: string;
  tipShow = false;
  leftSeconds: number = defaultLeftTime;

  constructor(private tipService: TipService) {
  }

  ngOnInit() {
    this.showTip();
  }

  showTip() {
    this.tipService.$tipQueue.subscribe(tip => {
      this.currentTip = tip;
      switch (tip.tipLevel) {
        case TipLevels.SUCCESS:
          this.currentTipLevel = 'success';
          break;
        case TipLevels.ERROR:
          this.currentTipLevel = 'error';
      }
      this.tipShow = true;
      const timer = setInterval(() => {
        this.leftSeconds--;
        if (this.leftSeconds < 0 || !this.tipService) {
          this.tipShow = false;
          clearInterval(timer);
          this.leftSeconds = defaultLeftTime;
        }
      }, defaultInterval);
    });
  }

  close() {
    this.tipShow = false;
  }

}
