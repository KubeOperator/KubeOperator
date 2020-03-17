import {Component, OnInit} from '@angular/core';
import {MessageCenterService} from './message-center.service';

@Component({
  selector: 'app-message-center',
  templateUrl: './message-center.component.html',
  styleUrls: ['./message-center.component.css']
})
export class MessageCenterComponent implements OnInit {

  constructor(private messageCenterService: MessageCenterService) {
  }


  ngOnInit() {
  }


}
