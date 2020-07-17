import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {Plan, PlanCreateRequest, PlanVmConfig, VmConfig} from '../plan';
import {PlanService} from '../plan.service';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {RegionService} from '../../region/region.service';
import {Region} from '../../region/region';
import {Zone} from '../../zone/zone';
import {NgForm} from '@angular/forms';
import {ClrWizard} from '@clr/angular';
import {ZoneService} from '../../zone/zone.service';

@Component({
    selector: 'app-plan-create',
    templateUrl: './plan-create.component.html',
    styleUrls: ['./plan-create.component.css']
})
export class PlanCreateComponent extends BaseModelComponent<Plan> implements OnInit {

    opened = false;
    regions: Region[] = [];
    zones: Zone[] = [];
    item: PlanCreateRequest = new PlanCreateRequest();
    vmConfigs: PlanVmConfig[] = [];
    regionName: string;
    @Output() created = new EventEmitter();
    @ViewChild('basicForm', {static: true}) basicForm: NgForm;
    @ViewChild('planForm', {static: true}) planForm: NgForm;
    @ViewChild('wizard') wizard: ClrWizard;


    constructor(private planService: PlanService, private modalAlertService: ModalAlertService, private regionService: RegionService,
                private translateService: TranslateService, private commonAlertService: CommonAlertService,
                private zoneService: ZoneService) {
        super(planService);
    }

    ngOnInit(): void {
    }

    open() {
        this.opened = true;
        this.regionService.list().subscribe(res => {
            this.regions = res.items;
        }, error => {

        });
    }

    onCancel() {
        this.opened = false;
        this.created.emit();
        this.wizard.reset();
        this.item = new PlanCreateRequest();
        this.basicForm.resetForm(this.item);
        this.planForm.resetForm(this.item);
    }

    onSubmit() {
        if (this.item.deployTemplate === 'SINGLE') {
            this.item.zones = [];
            this.item.zones.push(this.item.zone);
            this.onCancel();
        }
        this.planService.create(this.item).subscribe(res => {
            this.onCancel();
        });
    }

    onBasicFormCommit() {
        this.listZones();
        this.listVmConfigs();
    }

    onRegionChange() {
        this.regions.forEach(region => {
            if (region.id === this.item.regionId) {
                this.regionName = region.name;
            }
        });
    }

    onDeployChange() {

    }

    listVmConfigs() {
        this.planService.listVmConfigs(this.regionName).subscribe(res => {
            this.vmConfigs = res;
        });
    }

    listZones() {
        this.zoneService.listByRegionId(this.item.regionId).subscribe(res => {
            this.zones = res;
        });
    }

}
