import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {MessageCenterService} from '../../message-center.service';
import {CommonAlertService} from '../../../base/header/common-alert.service';
import {AlertLevels} from '../../../base/header/components/common-alert/alert';


@Component({
  selector: 'app-subscribe-config',
  templateUrl: './subscribe-config.component.html',
  styleUrls: ['./subscribe-config.component.css']
})
export class SubscribeConfigComponent implements OnInit {

  @Input() showConfigModal = false;
  @Input() subscribeConfig = {};
  @Output() subscribeConfigChange = new EventEmitter();


  constructor(private  messageCenterService: MessageCenterService, private alertService: CommonAlertService) {
  }

  ngOnInit() {
  }

  onCancel() {
    this.showConfigModal = false;
    this.subscribeConfigChange.emit(this.showConfigModal);
  }

  onSubmit(subscribable) {
    this.messageCenterService.updateSubscribe(subscribable).subscribe(res => {
      this.alertService.showAlert('更新成功', AlertLevels.SUCCESS);
      this.onCancel();
    });
  }

  changeValue(vars, type) {
    if (vars === 'DISABLE') {
      vars = 'ENABLE';
    } else {
      vars = 'DISABLE';
    }
    this.subscribeConfig['vars'][type] = vars;
  }
}
