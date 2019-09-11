import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {CloudTemplate, Region} from '../../region/region';
import {CloudTemplateService} from '../../region/cloud-template.service';
import {RegionService} from '../../region/region.service';
import {CloudService} from '../../region/cloud.service';
import {Plan} from '../plan';
import {ZoneService} from '../../zone/zone.service';
import {Zone} from '../../zone/zone';
import {PlanService} from '../plan.service';
import {NgForm} from '@angular/forms';
import {ClrWizard} from '@clr/angular';
import {catchError} from 'rxjs/operators';

@Component({
  selector: 'app-plan-create',
  templateUrl: './plan-create.component.html',
  styleUrls: ['./plan-create.component.css']
})
export class PlanCreateComponent implements OnInit {

  @Output() create = new EventEmitter<boolean>();
  createOpened: boolean;
  isSubmitGoing = false;
  item: Plan = new Plan();
  loading = false;
  cloudTemplate: CloudTemplate;
  regions: Region[] = [];
  region: Region;
  zones: Zone[] = [];
  zone: Zone;
  @ViewChild('basicForm') basicForm: NgForm;
  @ViewChild('planForm') planForm: NgForm;
  @ViewChild('wizard') wizard: ClrWizard;

  constructor(private cloudTemplateService: CloudTemplateService, private regionService: RegionService,
              private cloudService: CloudService, private zoneService: ZoneService, private planService: PlanService) {
  }

  ngOnInit() {
  }

  get nameCtrl() {
    return this.basicForm.controls['name'];
  }

  newItem() {
    this.basicForm.resetForm();
    this.wizard.reset();
    this.item = new Plan();
    this.regions = [];
    this.listRegion();
    this.createOpened = true;
    this.cloudTemplate = null;
  }

  listRegion() {
    this.regionService.listRegion().subscribe(data => {
      this.regions = data;
    });
  }

  nameOnBlur() {
    this.planService.getPlan(this.item.name).pipe(catchError(() => null)).subscribe((data) => {
      if (this.item.name) {
        this.nameCtrl.setErrors({repeat: true});
      }
    });
  }

  onRegionChange() {
    this.item.zone = undefined;
    this.item.zones = [];
    this.regions.forEach(region => {
      if (this.item.region === region.name) {
        this.region = region;
        this.cloudTemplateService.getCloudTemplate(region.template).subscribe(data => {
          this.cloudTemplate = data;
        });
      }
    });
  }

  onDeployChange() {
    this.item.zone = undefined;
    this.item.zones = [];
  }

  onBasicFormCommit() {
    this.zoneService.listZones().subscribe(data => {
      this.zones = data.filter(zone => {
        return zone.region === this.region.name && zone.status === 'READY';
      });
    });
  }

  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.isSubmitGoing = true;
    this.planService.createPlan(this.item).subscribe(data => {
      this.isSubmitGoing = false;
      this.createOpened = false;
      this.create.emit(true);
    });
  }

  onCancel() {
    this.createOpened = false;
  }

}
