import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Host} from '../host';
import {HostService} from '../host.service';
import {TipService} from '../../tip/tip.service';
import {TipLevels} from '../../tip/tipLevels';

@Component({
  selector: 'app-host-create',
  templateUrl: './host-create.component.html',
  styleUrls: ['./host-create.component.css']
})
export class HostCreateComponent implements OnInit {

  constructor(private hostService: HostService, private tipService: TipService) {
  }

  @Output() create = new EventEmitter<boolean>();
  staticBackdrop = true;
  closable = false;
  createHostOpened: boolean;
  isSubmitGoing = false;
  host: Host = new Host();
  loading = false;

  ngOnInit() {
  }


  onCancel() {
    this.createHostOpened = false;
  }

  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.isSubmitGoing = true;
    this.loading = true;
    this.hostService.createHost(this.host).subscribe(data => {
      this.createHostOpened = false;
      this.isSubmitGoing = false;
      this.create.emit(true);
      this.loading = false;
      this.tipService.showTip('创建主机成功', TipLevels.SUCCESS);
    }, err => {
      this.createHostOpened = false;
      this.isSubmitGoing = false;
      this.create.emit(true);
      this.loading = false;
      this.tipService.showTip('创建主机失败:' + err, TipLevels.ERROR);
    });
  }

  newHost() {
    this.host = new Host();
    this.createHostOpened = true;
  }
}
