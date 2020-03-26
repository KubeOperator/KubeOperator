import {Component, OnInit, ViewChild} from '@angular/core';
import {MessageCenterService} from '../message-center.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {LocalMailDetailComponent} from "./local-mail-detail/local-mail-detail.component";

@Component({
  selector: 'app-local-mail',
  templateUrl: './local-mail.component.html',
  styleUrls: ['./local-mail.component.css']
})
export class LocalMailComponent implements OnInit {

  loading = false;
  messages = [];
  selectedMessages = [];
  limit = 10;
  page = 1;
  type = 'ALL';
  readStatus = 'ALL';
  level = 'ALL';
  total = 0;

  @ViewChild(LocalMailDetailComponent, {static: true})
  detail: LocalMailDetailComponent;

  constructor(private messageCenterService: MessageCenterService, private alertService: CommonAlertService) {
  }

  ngOnInit() {
    this.listMessage(this.limit, this.type, this.readStatus, this.level);
  }

  listMessage(limit, type, readStatus, level) {
    this.limit = limit;
    this.type = type;
    this.readStatus = readStatus;
    this.level = level;
    this.loading = true;
    this.messageCenterService.listUserMessageByPage(this.limit, this.page, this.type, this.readStatus, this.level).subscribe(res => {
      this.messages = res.data;
      this.total = res.total;
      // this.page = res.page_num;
      this.loading = false;
    });
  }

  updateMessage() {
    const promises: Promise<{}>[] = [];
    this.selectedMessages.forEach(msg => {
      promises.push(this.messageCenterService.updateUserMessage(msg).toPromise());
    });
    Promise.all(promises).then(() => {
      this.listMessage(this.limit, this.type, this.readStatus, this.level);
      this.alertService.showAlert('更新成功', AlertLevels.SUCCESS);
    }, res => {
      this.alertService.showAlert('更新失败' + res.error.msg, AlertLevels.ERROR);
    }).finally(() => {
      this.selectedMessages = [];
    });
  }

  updateAllMessage() {
    this.messageCenterService.updateAllUserMessage().subscribe(data => {
      this.alertService.showAlert('更新成功', AlertLevels.SUCCESS);
      this.listMessage(this.limit, this.type, this.readStatus, this.level);
    });
  }

  showDetail(message) {
    this.updateSingleMessage(message);
    const detailMessage = JSON.parse(JSON.stringify(message));
    this.detail.message = detailMessage;
    this.detail.message.message_detail.content = JSON.parse(detailMessage.message_detail.content);
    this.detail.message.message_detail.content.detail = JSON.parse(this.detail.message.message_detail.content.detail);
    this.detail.open = true;
  }

  updateSingleMessage(message) {
    message.read_status = 'READ';
    this.messageCenterService.updateUserMessage(message).subscribe(data => {
    });
  }

  getTitle(message) {
    const detailMessage = JSON.parse(JSON.stringify(message));
    const content = JSON.parse(detailMessage.message_detail.content);
    let title = message.message_detail.title;
    if (content.resource_type === 'CLUSTER' || content.resource_type === 'CLUSTER_EVENT' || content.resource_type === 'CLUSTER_USAGE') {
      title = title + '[集群:' + content.resource_name + ']';
    }
    return title;
  }
}
