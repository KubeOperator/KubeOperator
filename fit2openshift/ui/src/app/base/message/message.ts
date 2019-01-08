import {MessageLevels} from './message-level';

export class Message {
  msg: string;
  level: MessageLevels;

  constructor(msg: string, level: MessageLevels) {
    this.msg = msg;
    this.level = level;
  }
}
