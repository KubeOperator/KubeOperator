import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {HostInfo} from '../host';
import {HostInfoService} from './host-info.service';
import {HostService} from '../host.service';

@Component({
  selector: 'app-host-info',
  templateUrl: './host-info.component.html',
  styleUrls: ['./host-info.component.css']
})
export class HostInfoComponent implements OnInit {

  hostId: string;
  hostInfo: HostInfo = null;
  loading = false;
  errorText = null;
  @Input() showInfoModal = false;
  @Output() showInfoModalChange = new EventEmitter();

  constructor(private hostInfoService: HostInfoService, private hostService: HostService) {
  }

  ngOnInit() {
  }

  loadHostInfo() {
    this.hostService.getHost(this.hostId).subscribe(data => {
      this.hostInfo = data.info;
    });


  }

  update() {
    this.loading = true;
    this.hostInfoService.loadHostInfo(this.hostId).subscribe(data => {
      this.loading = false;
      this.hostInfo = data;
    }, error => {
      this.errorText = error;
    });
  }

  cancel() {
    this.showInfoModal = false;
    this.showInfoModalChange.emit(this.showInfoModal);
  }


}
