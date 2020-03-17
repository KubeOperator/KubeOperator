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
  showConfig = false;
  subscribeConfig = {};

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

  openModal(subscribe) {
    this.showConfig = true;
    this.subscribeConfig = subscribe;
  }

  getData(showConfig) {
    this.showConfig = showConfig;
    this.listSubscribe();
  }
}
