import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Region} from '../../region/region';
import {NgForm} from '@angular/forms';
import {ClrWizard} from '@clr/angular';
import {RegionService} from '../../region/region.service';
import {CloudService} from '../../region/cloud.service';
import {Zone} from '../zone';
import {CloudZone, Subnet} from '../../region/cloud';
import {ZoneService} from '../zone.service';
import {catchError} from 'rxjs/operators';
import * as ipaddr from 'ipaddr.js';
import {IpService} from '../ip.service';
import * as globals from '../../globals';


@Component({
  selector: 'app-zone-create',
  templateUrl: './zone-create.component.html',
  styleUrls: ['./zone-create.component.css']
})
export class ZoneCreateComponent implements OnInit {

  @Output() create = new EventEmitter<boolean>();
  createOpened: boolean;
  isSubmitGoing = false;
  item: Zone = new Zone();
  cloudZones: CloudZone[] = [];
  cloudZone: CloudZone;
  regions: Region[] = [];
  region: Region = new Region();
  subnetList: Subnet[] = [];
  loading = false;
  networkValid = false;
  networkError = [];
  @ViewChild('basicForm', {static: true}) basicForm: NgForm;
  @ViewChild('wizard', {static: true}) wizard: ClrWizard;
  name_pattern = globals.host_name_pattern;
  name_pattern_tip = globals.host_name_pattern_tip;

  constructor(private regionService: RegionService,
              private cloudService: CloudService,
              private zoneService: ZoneService,
              private ipService: IpService) {
  }

  ngOnInit() {
  }


  get nameCtrl() {
    return this.basicForm.controls['name'];
  }

  nameOnBlur() {
    this.zoneService.getZone(this.item.name).pipe(catchError(() => null)).subscribe((data) => {
      if (this.item.name) {
        this.nameCtrl.setErrors({repeat: true});
      }
    });
  }

  newItem() {
    this.item = new Zone();
    this.reset();
    this.createOpened = true;
    this.listRegion();
  }

  reset() {
    this.wizard.reset();
    this.basicForm.resetForm();
    this.regions = [];
    this.cloudZones = [];
    this.cloudZone = null;
    this.subnetList = [];
  }

  checkNetwork() {
    this.networkError = [];
    let result = true;
    const ipStart = this.item.vars['ip_start'];
    const ipEnd = this.item.vars['ip_end'];
    if (!ipaddr.isValid(ipStart)) {
      result = false;
      this.networkError.push('起始IP不是有效的IPV4地址!');
    }
    if (!ipaddr.isValid(ipEnd)) {
      result = false;
      this.networkError.push('截止IP不是有效的IPV4地址!');
    }
    if (ipaddr.isValid(ipStart) && ipaddr.isValid(ipEnd)) {
      const start = ipaddr.parse(ipStart);
      const end = ipaddr.parse(ipEnd);
      if (start >= end) {
        result = false;
        this.networkError.push('截止IP必须大于起始IP!');
      }
    }
    if (this.region.template === 'vsphere') {
      const mask = this.item.vars['net_mask'];
      const gateway = this.item.vars['vc_gateway'];
      if (!ipaddr.isValid(gateway)) {
        result = false;
        this.networkError.push('网关IP不是有效的IPV4地址!');
      }
      if (!ipaddr.isValid(mask)) {
        result = false;
        this.networkError.push('子网掩码不是有效的IPV4地址!');
      } else {
        const maskIp = ipaddr.parse(mask);
        if (maskIp.prefixLengthFromSubnetMask() == null) {
          result = false;
          this.networkError.push('子网掩码无效！');
        }
      }
    }
    this.networkValid = result;
  }

  listRegion() {
    this.regionService.listRegion().subscribe(data => {
      this.regions = data;
    });
  }

  onRegionChange() {
    this.regions.forEach(region => {
      if (region.name === this.item.region) {
        this.region = region;
      }
    });
  }

  onComputeChange() {
    this.item.vars = {};
    this.cloudZones.forEach(zone => {
      if (this.item.cluster === zone.cluster) {
        this.cloudZone = zone;
      }
    });
  }

  onNetworkChange() {
    this.cloudZone.networkList.forEach(network => {
      if (this.item.vars['openstack_network'] === network.id) {
        this.subnetList = network.subnetList;
      }
    });
  }

  onBasicFormCommit() {
    this.loading = true;
    this.cloudService.listZone(this.item.region).subscribe(data => {
      this.loading = false;
      this.cloudZones = data;
    });
  }

  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.isSubmitGoing = true;
    switch (this.region.template) {
      case 'vsphere':
        this.item.vars['vc_cluster'] = this.item.cluster;
        break;
      case 'openstack':
        this.item.vars['openstack_zone'] = this.item.cluster;
        break;
    }
    this.zoneService.createZones(this.item).subscribe(data => {
      this.isSubmitGoing = false;
      this.createOpened = false;
      this.create.emit(true);
    });
  }

  onCancel() {
    this.createOpened = false;
  }

}
