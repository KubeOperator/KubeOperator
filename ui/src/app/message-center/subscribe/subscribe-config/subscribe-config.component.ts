import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {MessageCenterService} from '../../message-center.service';

@Component({
  selector: 'app-subscribe-config',
  templateUrl: './subscribe-config.component.html',
  styleUrls: ['./subscribe-config.component.css']
})
export class SubscribeConfigComponent implements OnInit {

  @Input() showConfigModal = false;
  @Input() subscribeConfig = {};
  @Output() subscribeConfigChange = new EventEmitter();


  constructor(private  messageCenterService: MessageCenterService) {
  }

  ngOnInit() {
  }

  onCancel() {
    this.showConfigModal = false;
    this.subscribeConfigChange.emit(this.showConfigModal);
  }

  onSubmit(subscribable) {
    this.messageCenterService.updateSubscribe(subscribable).subscribe(res => {
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
