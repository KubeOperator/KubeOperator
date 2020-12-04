import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
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
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {ProjectService} from '../../../project/project.service';
import {Project} from '../../../project/project';
import {NamePattern} from '../../../../constant/pattern';

@Component({
    selector: 'app-plan-create',
    templateUrl: './plan-create.component.html',
    styleUrls: ['./plan-create.component.css']
})
export class PlanCreateComponent extends BaseModelDirective<Plan> implements OnInit {

    namePattern = NamePattern;
    opened = false;
    regions: Region[] = [];
    zones: Zone[] = [];
    item: PlanCreateRequest = new PlanCreateRequest();
    vmConfigs: PlanVmConfig[] = [];
    regionName: string;
    projects: Project[] = [];
    regionId: string;
    currentProvider: string;
    isSubmitGoing = false;
    @Output() created = new EventEmitter();
    @ViewChild('basicForm', {static: true}) basicForm: NgForm;
    @ViewChild('planForm', {static: true}) planForm: NgForm;
    @ViewChild('wizard') wizard: ClrWizard;


    constructor(private planService: PlanService, private modalAlertService: ModalAlertService, private regionService: RegionService,
                private translateService: TranslateService, private commonAlertService: CommonAlertService,
                private zoneService: ZoneService, private projectService: ProjectService) {
        super(planService);
    }

    ngOnInit(): void {
    }

    open() {
        this.opened = true;
        this.listProjects();
        this.listRegions();
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
        }
        this.isSubmitGoing = true;
        this.planService.create(this.item).subscribe(res => {
            this.onCancel();
            this.created.emit();
            this.isSubmitGoing = false;
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    onBasicFormCommit() {
        this.listZones();
        this.listVmConfigs();
    }

    onRegionChange() {
        this.regions.forEach(region => {
            if (region.name === this.item.region) {
                this.regionName = region.name;
                this.regionId = region.id;
                this.currentProvider = region.provider;
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
        this.zoneService.listByRegionName(this.regionName).subscribe(res => {
            this.zones = res;
        });
    }

    listProjects() {
        this.projectService.list().subscribe(res => {
            this.projects = res.items;
        }, error => {

        });
    }

    listRegions() {
        this.regionService.list().subscribe(res => {
            this.regions = res.items;
        }, error => {

        });
    }
}
