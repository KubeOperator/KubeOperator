import {Injectable} from '@angular/core';
import {Subject} from 'rxjs';
import {Tip} from './tip';
import {TipLevels} from './tipLevels';

@Injectable({
  providedIn: 'root'
})
export class TipService {

  tipQueue = new Subject<Tip>();
  $tipQueue = this.tipQueue.asObservable();

  constructor() {
  }

  showTip(msg: string, level: TipLevels) {
    this.tipQueue.next(new Tip(msg, level));
  }

}
