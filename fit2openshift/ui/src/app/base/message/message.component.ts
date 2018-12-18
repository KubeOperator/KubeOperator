import {Component, OnInit} from '@angular/core';
import {Message} from './message';
import {MessageService} from '../message.service';
import {MessageLevels} from './message-level';

@Component({
  selector: 'app-message',
  templateUrl: './message.component.html',
  styleUrls: ['./message.component.css']
})
export class MessageComponent implements OnInit {
  messageShow = false;
  currentMessage: Message;

  constructor(private msgService: MessageService) {
  }

  ngOnInit() {
    this.showMassage();
  }

  showMassage() {
    this.msgService.$messageQueue.subscribe(msg => {
      this.currentMessage = msg;
      this.messageShow = true;
    });
  }

  close() {
    this.messageShow = false;
  }

}
