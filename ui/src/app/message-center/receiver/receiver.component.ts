import {Component, OnInit} from '@angular/core';
import {MessageCenterService} from '../message-center.service';

@Component({
  selector: 'app-receiver',
  templateUrl: './receiver.component.html',
  styleUrls: ['./receiver.component.css']
})
export class ReceiverComponent implements OnInit {

  userConfig = {};
  submitGoing = false;


  constructor(private messageCenterService: MessageCenterService) {
  }

  ngOnInit() {
    this.listReceiver();
  }

  onCancel() {

  }

  onSubmit() {
    this.messageCenterService.updateUserReceiver(this.userConfig).subscribe(res => {

    });
  }

  listReceiver() {
    this.messageCenterService.listUserReceiver().subscribe(data => {
      this.userConfig = data[0];
    });
  }

}
