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
    @Output() created = new EventEmitter();
    @ViewChild('basicForm', {static: true}) basicForm: NgForm;
    @ViewChild('wizard') wizard: ClrWizard;


    constructor(private planService: PlanService, private modalAlertService: ModalAlertService, private regionService: RegionService,
                private translateService: TranslateService, private commonAlertService: CommonAlertService,
                private zoneService: ZoneService) {
        super(planService);
    }

    ngOnInit(): void {
    }

    open() {
        this.wizard.reset();
        this.opened = true;
        this.regionService.list().subscribe(res => {
            this.regions = res.items;
        }, error => {

        });
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        if (this.item.deployTemplate === 'SINGLE') {
            this.item.zones = [];
            this.item.zones.push(this.item.zone);
        }
        this.planService.create(this.item).subscribe(res => {
            this.opened = false;
        });
    }

    onBasicFormCommit() {
        this.listZones();
        this.listVmConfigs();
    }

    onRegionChange() {

    }

    onDeployChange() {

    }

    listVmConfigs() {
        this.planService.listVmConfigs().subscribe(res => {
            this.vmConfigs = res;
        });
    }

    listZones() {
        this.zoneService.list().subscribe(res => {
            this.zones = res.items;
        });
    }

}
