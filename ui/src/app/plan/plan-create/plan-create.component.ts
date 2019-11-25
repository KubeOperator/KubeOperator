import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Region} from '../../region/region';
import {RegionService} from '../../region/region.service';
import {CloudService} from '../../region/cloud.service';
import {ComputeModel, Plan} from '../plan';
import {ZoneService} from '../../zone/zone.service';
import {Zone} from '../../zone/zone';
import {PlanService} from '../plan.service';
import {NgForm} from '@angular/forms';
import {ClrWizard} from '@clr/angular';
import {catchError} from 'rxjs/operators';
import * as globals from '../../globals';


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
  regions: Region[] = [];
  region: Region;
  zones: Zone[] = [];
  computeModels: ComputeModel[] = [];
  zone: Zone;
  @ViewChild('basicForm', {static: true}) basicForm: NgForm;
  @ViewChild('planForm', {static: true}) planForm: NgForm;
  @ViewChild('wizard', {static: true}) wizard: ClrWizard;
  name_pattern = globals.host_name_pattern;
  name_pattern_tip = globals.host_name_pattern_tip;

  constructor(private regionService: RegionService,
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
    this.listComputeModel();
    this.createOpened = true;
  }

  listRegion() {
    this.regionService.listRegion().subscribe(data => {
      this.regions = data;
    });
  }

  listComputeModel() {
    this.planService.getComputeModel().subscribe(data => {
      console.log(data);
      this.computeModels = data;
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
        if (this.region.template === 'openstack') {
          this.setFlavorModels();
        }
      }
    });
  }

  setFlavorModels() {
    this.computeModels = [];
    this.cloudService.listFlavor(this.region.name).subscribe(data => {
      this.computeModels = data;
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
    if (this.region.template === 'openstack') {
      this.item.vars['compute_models'] = this.computeModels;
    }
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
