import {Component, Input, OnInit, Output} from '@angular/core';
import {PackageService} from '../../package/package.service';
import {Package} from '../../package/package';
import {ConfigBase} from './config-base';
import {FormControl, FormGroup} from '@angular/forms';
import {ConfigControlService} from './config-control.service';
import {TextConfig} from './config-text';
import {SelectConfig} from './config-select';

@Component({
  selector: 'app-cluster-config',
  templateUrl: './cluster-config.component.html',
  styleUrls: ['./cluster-config.component.css']
})
export class ClusterConfigComponent implements OnInit {
  offlines: Package[] = [];
  configs: ConfigBase<any>[] = [];
  form = new FormGroup({
    offline: new FormControl()
  });


  constructor(private offlineService: PackageService, private ccs: ConfigControlService) {
  }

  ngOnInit() {
    this.listOffline();
  }

  listOffline() {
    this.offlineService.listPackage().subscribe(data => this.offlines = data);
  }

  onSubmit() {
    console.log(this.form.value);
  }


}

