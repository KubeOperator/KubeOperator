import {TipLevels} from './tipLevels';

export class Tip {
  msg: string;
  tipLevel: TipLevels;


  constructor(msg: string, tipLevel: TipLevels) {
    this.msg = msg;
    this.tipLevel = tipLevel;
  }
}
