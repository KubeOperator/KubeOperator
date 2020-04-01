import {Component, OnInit} from '@angular/core';
import {MessageCenterService} from '../message-center.service';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';


@Component({
  selector: 'app-receiver',
  templateUrl: './receiver.component.html',
  styleUrls: ['./receiver.component.css']
})
export class ReceiverComponent implements OnInit {

  userConfig = {};
  submitGoing = false;


  constructor(private messageCenterService: MessageCenterService, private alertService: CommonAlertService) {
  }

  ngOnInit() {
    this.listReceiver();
  }

  onCancel() {

  }

  onSubmit() {
    this.messageCenterService.updateUserReceiver(this.userConfig).subscribe(res => {
      this.alertService.showAlert('更新成功', AlertLevels.SUCCESS);
    });
  }

  listReceiver() {
    this.messageCenterService.listUserReceiver().subscribe(data => {
      this.userConfig = data[0];
    });
  }

}
