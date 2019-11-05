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

  constructor(private hostService: HostService) {
  }

  ngOnInit() {
  }

  refresh() {
    this.loading = true;
    this.hostService.getHost(this.host.id).subscribe(data => {
      this.loading = false;
      this.host = data;
    });
  }

  cancel() {
    this.showInfoModal = false;
    this.showInfoModalChange.emit(this.showInfoModal);
  }


}
