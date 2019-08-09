import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {CloudTemplate, Region} from '../../region/region';
import {NgForm} from '@angular/forms';
import {ClrWizard} from '@clr/angular';
import {RegionService} from '../../region/region.service';
import {CloudService} from '../../region/cloud.service';
import {Zone} from '../zone';
import {CloudZone} from '../../region/cloud';
import {CloudTemplateService} from '../../region/cloud-template.service';
import {ZoneService} from '../zone.service';
import {catchError} from 'rxjs/operators';

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
  region: Region;
  cloudTemplate: CloudTemplate;
  loading = false;
  @ViewChild('basicForm') basicForm: NgForm;
  @ViewChild('wizard') wizard: ClrWizard;

  constructor(private regionService: RegionService,
              private cloudService: CloudService,
              private cloudTemplateService: CloudTemplateService,
              private zoneService: ZoneService) {
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
        this.cloudTemplateService.getCloudTemplate(region.template).subscribe(data => {
          this.cloudTemplate = data;
        });
      }
    });
  }

  onComputeChange() {
    this.cloudZones.forEach(zone => {
      if (this.item.cloud_zone === zone.name) {
        this.cloudZone = zone;
        this.cloudTemplate.meta.zone.network.configs.forEach(config => {
          if (config.option_var === 'networks') {
            config.objs = zone.networks;
          }
        });
        this.cloudTemplate.meta.zone.storage.configs.forEach(config => {
          if (config.option_var === 'storages') {
            config.objs = zone.storages;
          }
        });
        this.cloudTemplate.meta.zone.image.configs.forEach(config => {
          if (config.option_var === 'images') {
            config.objs = zone.images;
          }
        });
        this.item.vars = {};
      }
    });
  }

  onBasicFormCommit() {
    this.cloudService.listZone(this.item.region).subscribe(data => {
      this.cloudZones = data;
    });
  }

  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.isSubmitGoing = true;
    this.zoneService.createZones(this.item).subscribe(data => {
      this.isSubmitGoing = false;
      this.createOpened = false;
      this.create.emit(true);
    });
  }

  onCancel() {
    this.createOpened = false;
    this.reset();
  }

}
