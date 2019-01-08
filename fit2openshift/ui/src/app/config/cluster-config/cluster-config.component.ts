import {Component, Input, OnInit, Output} from '@angular/core';
import {OfflineService} from '../../offline/offline.service';
import {Offline} from '../../offline/Offline';
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
  offlines: Offline[] = [];
  configs: ConfigBase<any>[] = [];
  form = new FormGroup({
    offline: new FormControl()
  });


  constructor(private offlineService: OfflineService, private ccs: ConfigControlService) {
  }

  ngOnInit() {
    this.listOffline();
    this.getConfig();
  }

  listOffline() {
    this.offlineService.listOfflines().subscribe(data => this.offlines = data);
  }

  onSubmit() {
    console.log(this.form.value);
  }

  getConfig() {
    this.form.get('offline').valueChanges.subscribe(data => {
      this.offlineService.getOffline('aa').subscribe(data => {
        this.configs = data.config;
        this.form.addControl('config', this.ccs.toFormGroup(this.configs));
      });
    });
  }

}

