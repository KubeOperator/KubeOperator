import {Component, OnInit} from '@angular/core';
import {Dns} from './dns';
import {DnsService} from './dns.service';
import {CommonAlertService} from '../base/header/common-alert.service';
import {AlertLevels} from '../base/header/components/common-alert/alert';

@Component({
  selector: 'app-dns',
  templateUrl: './dns.component.html',
  styleUrls: ['./dns.component.css']
})
export class DnsComponent implements OnInit {

  constructor(private dnsService: DnsService, private alert: CommonAlertService) {
  }

  dns: Dns = new Dns();

  ngOnInit() {
    this.reset();
  }

  reset() {
    this.dnsService.getDns().subscribe(data => {
      this.dns = data;
    });
  }

  onSubmit() {
    this.dnsService.updateDns(this.dns).subscribe(data => {
      this.alert.showAlert('修改成功！', AlertLevels.SUCCESS);
      this.reset();
    });
  }

}
