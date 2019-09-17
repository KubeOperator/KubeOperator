import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {HostInfoService} from './host-info.service';
import {HostService} from '../host.service';
import {Host} from '../host';

@Component({
  selector: 'app-host-info',
  templateUrl: './host-info.component.html',
  styleUrls: ['./host-info.component.css']
})
export class HostInfoComponent implements OnInit {

  host: Host = new Host;
  loading = false;
  @Input() showInfoModal = false;
  @Output() showInfoModalChange = new EventEmitter();

  constructor() {
  }

  ngOnInit() {
  }

  cancel() {
    this.showInfoModal = false;
    this.showInfoModalChange.emit(this.showInfoModal);
  }


}
