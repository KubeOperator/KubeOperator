import {Injectable} from '@angular/core';
import {Subject} from 'rxjs';
import {Message} from './message/message';
import {MessageLevels} from './message/message-level';

@Injectable({
  providedIn: 'root'
})
export class MessageService {

  messagesQueue = new Subject<Message>();
  $messageQueue = this.messagesQueue.asObservable();

  constructor() {
  }

  announceMessage(msg: string, level: MessageLevels) {
    this.messagesQueue.next(new Message(msg, level));
  }
}
