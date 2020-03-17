import {Component, OnInit} from '@angular/core';
import {MessageCenterService} from '../message-center.service';

@Component({
  selector: 'app-subscribe',
  templateUrl: './subscribe.component.html',
  styleUrls: ['./subscribe.component.css']
})
export class SubscribeComponent implements OnInit {

  loading = false;
  subscribes = [];


  constructor(private messageCenterService: MessageCenterService) {
  }


  ngOnInit() {
    this.listSubscribe();
  }

  listSubscribe() {
    this.loading = true;
    this.messageCenterService.listSubscribe().subscribe(data => {
      this.subscribes = data;
      this.loading = false;
    });
  }
}
